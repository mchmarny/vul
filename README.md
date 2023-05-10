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

## Schedule 

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


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
