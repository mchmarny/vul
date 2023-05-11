# vul

## setup

Create a secret that will hold the snyk token:

```shell
gcloud secrets create vimp-snyk-token --replication-policy="automatic"
echo -n "${SNYK_TOKEN}" | gcloud secrets versions add vimp-snyk-token --data-file=- \
export SNYK_SECRET=$(gcloud secrets versions describe 1 --secret vimp-snyk-token --format="value(NAME)")
```

Create pub/sub topic which will be used to queue your images: 

> The topic name can be anything as long as it's the same in the topic and trigger create commands.

```shell
gcloud pubsub topics create image-queue --project $PROJECT_ID
```

Next create a trigger to process any new events on that queue with the same build config as above: 

```shell
gcloud alpha builds triggers create pubsub \
    --project=$PROJECT_ID \
    --region=$REGION \
    --name=scan-and-save-external-image-exposure-data \
    --topic=projects/$PROJECT_ID/topics/image-queue \
    --build-config=config/process.yaml \
    --substitutions=_DIGEST='$(body.message.data)' \
    --repo=https://www.github.com/mchmarny/vul \
    --repo-type=GITHUB \
    --branch=main
```

Now to process new image simply publish the image URI (with digest) to that topic:

```shell
gcloud pubsub topics publish image-queue \
    --message=https://docker.io/redis \
    --project=$PROJECT_ID
```

## schedule 

To schedule [config/image.txt](config/image.txt) to be queued for processing: 

```shell
gcloud beta builds triggers create manual \
    --name=queue-images \
    --project=$PROJECT_ID \
    --region=$REGION \
    --repo=https://www.github.com/mchmarny/vul \
    --build-config=config/queue.yaml \
    --repo-type=GITHUB \
    --branch=main
```

Next, capture the trigger ID:

```shell
export TRIGGER_ID=$(gcloud beta builds triggers describe queue-images \
    --project=$PROJECT_ID --region=$REGION --format='value(id)')
```

You can run this trigger now manually, by invoking from `curl`. 

> This assumes that you have the necessary role to execute the build.

```shell
curl -X POST -H "Authorization: Bearer $(gcloud auth print-access-token)" \
     "https://cloudbuild.googleapis.com/v1/projects/$PROJECT_ID/locations/$REGION/triggers/$TRIGGER_ID:run"
```

That means we can now set it up as a Cloud Schedule, first, make sure the Cloud Build account has sufficient rights to execute the job:


```shell
export PROJECT_NUMBER=$(gcloud projects describe $PROJECT_ID --format='value(projectNumber)')
gcloud projects add-iam-policy-binding $PROJECT_ID \
    --member="serviceAccount:$PROJECT_NUMBER-compute@developer.gserviceaccount.com" \
    --role="roles/cloudbuild.builds.editor" \
    --condition=None
```

Finally, create the Cloud Scheduler job:

```shell
gcloud scheduler jobs create http queue-images-schedule \
    --http-method POST \
    --schedule='0 3,15 * * *' \
    --location=$REGION \
    --uri=https://cloudbuild.googleapis.com/v1/projects/$PROJECT_ID/locations/$REGION/triggers/$TRIGGER_ID:run \
    --oauth-service-account-email=$PROJECT_NUMBER-compute@developer.gserviceaccount.com \
    --oauth-token-scope=https://www.googleapis.com/auth/cloud-platform
```

Now everyday, at 3am and 3pm UTC, the image will be rebuilt and the Cloud Workstation configuration updated with the latest image.

## query 

List images:

```shell
curl -X GET https://vul.thingz.io/api/v1/images
```

That returns: 

```json
{
    "version": "v0.1.0",
    "created": "2023-05-10T17:44:22.128263373Z",
    "data": [
        {
            "image": "docker.io/amazon/aws-cli",
            "version_count": 1,
            "first_reading": "2023-05-08T22:03:34.68755Z",
            "last_reading": "2023-05-10T12:45:40.134534Z"
        },
        {
            "image": "docker.io/amazon/aws-lambda-python",
            "version_count": 1,
            "first_reading": "2023-05-08T22:03:43.964732Z",
            "last_reading": "2023-05-10T12:45:43.322162Z"
        },
        ...
    ]
}
```

Then, using one of the returned versions: 

```shell
curl -s -H "Content-Type: application/json" \
     -d '{ "image": "docker.io/bitnami/mariadb" }' \
    https://vul.thingz.io/api/v1/versions
```

The result: 

```json
{
    "version": "v0.1.0",
    "created": "2023-05-10T17:53:28.910083767Z",
    "criteria": {
        "image": "docker.io/bitnami/mariadb"
    },
    "data": [
        {
            "digest": "sha256:19a6c75aa7efeaa833e40bb6aa8659d04e030299a5b11e2db9345de752599db3",
            "source_count": 3,
            "first_reading": "2023-05-09T22:03:20.943867Z",
            "last_reading": "2023-05-10T12:46:11.244266Z",
            "package_count": 69
        },
        {
            "digest": "sha256:97b0be98b4714e81dac9ac55513f4f87c627d88da09d90c708229835124a8215",
            "source_count": 3,
            "first_reading": "2023-05-08T22:03:32.725514Z",
            "last_reading": "2023-05-08T22:04:18.365187Z",
            "package_count": 69
        },
        ...
    ]
}
```

Next, using image digest return exposures: 

```shell
curl -s -H "Content-Type: application/json" \
     -d '{ "image": "docker.io/bitnami/mongodb", "digest": "sha256:419f129df0140834d89c94b29700c91f38407182137be480a0d6c6cbe2e0d00a" }' \
    https://vul.thingz.io/api/v1/exposures
```

The result: 

```json
{
    "version": "v0.1.7",
    "created": "2023-05-11T11:53:10.535212577Z",
    "criteria": {
        "image": "docker.io/bitnami/mongodb",
        "digest": "sha256:419f129df0140834d89c94b29700c91f38407182137be480a0d6c6cbe2e0d00a"
    },
    "data": {
        "CVE-2005-2541": [
            {
                "source": "grype",
                "package": "tar",
                "version": "1.34+dfsg-1",
                "severity": "high",
                "score": 10,
                "fixed": false
            },
            ...
    ],
}
```

You can also get a exposure timeline for each image: 

```shell
curl -s -H "Content-Type: application/json" \
     -d '{ "image": "docker.io/bitnami/mongodb" }' \
    https://vul.thingz.io/api/v1/timeline
```

The result: 

```json
{
    "version": "v0.1.6",
    "created": "2023-05-11T10:58:34.090456796Z",
    "criteria": {
        "image": "docker.io/bitnami/mongodb",
        "from_day": "2023-04-11",
        "to_day": "2023-05-11"
    },
    "data": {
        "2023-05-08": {
            "sources": {
                "grype": {
                    "total": 96,
                    "negligible": 1,
                    "low": 10,
                    "medium": 39,
                    "high": 40,
                    "critical": 5,
                    "unknown": 1
                },
                ...
                }
            }
        },
        ...
    }
}
```

## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
