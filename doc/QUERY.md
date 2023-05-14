# query 

## summary 

Get summary of all the data:

```shell
curl https://vul.thingz.io/api/v1/summary
```

Returns: 

```json
{
    "version": "v0.4.1",
    "created": "2023-05-14T12:38:54.796657223Z",
    "data": {
        "image_count": 17,
        "version_count": 39,
        "source_count": 3,
        "package_count": 576,
        "exposure": {
            "total": 18543,
            "negligible": 38,
            "low": 4659,
            "medium": 6699,
            "high": 6054,
            "critical": 988,
            "unknown": 105
        },
        "first_reading": "2023-05-08T15:43:02.546802Z",
        "last_reading": "2023-05-14T10:20:12.343577Z"
    }
}
```

## images 

List images currently being tracked:

```shell
curl https://vul.thingz.io/api/v1/images
```

Returns: 

```json
{
    "version": "v0.4.1",
    "created": "2023-05-14T12:36:23.046905103Z",
    "data": [
        "docker.io/dynatrace/oneagent",
        "docker.io/hashicorp/terraform",
        "docker.io/openjdk",
        ...
    ]
}
```

## image summary 

Get single image summary:

```shell
curl --data-urlencode "img=docker.io/openjdk" --get https://vul.thingz.io/api/v1/summary
```

Returns: 

```json
{
    "version": "v0.4.1",
    "created": "2023-05-14T12:41:56.358712828Z",
    "data": {
        "image": "docker.io/openjdk",
        "image_count": 1,
        "version_count": 1,
        "source_count": 3,
        "package_count": 7,
        "exposure": {
            "total": 43,
            "negligible": 0,
            "low": 0,
            "medium": 22,
            "high": 18,
            "critical": 3,
            "unknown": 0
        },
        "first_reading": "2023-05-08T22:03:51.573467Z",
        "last_reading": "2023-05-14T10:04:51.299231Z"
    }
}
```

## image timeline 

Get timeline of image exposures:

```shell
curl --data-urlencode "img=docker.io/openjdk" --get https://vul.thingz.io/api/v1/timeline
```

Returns: 

```json
{
    "version": "v0.4.1",
    "created": "2023-05-14T12:47:39.987826145Z",
    "criteria": {
        "image": "docker.io/openjdk",
        "since": "2023-04-14"
    },
    "data": [
        {
            "date": "2023-05-08",
            "grype": 9,
            "trivy": 14,
            "snyk": 20
        },
        {
            "date": "2023-05-09",
            "grype": 9,
            "trivy": 15,
            "snyk": 20
        },
        ...
    ]
}
```

## image versions

Get all versions for an image:

```shell
curl --data-urlencode "img=docker.io/openjdk" --get https://vul.thingz.io/api/v1/versions
```

Returns: 

```json
{
    "version": "v0.4.1",
    "created": "2023-05-14T12:42:59.676839242Z",
    "criteria": {
        "image": "docker.io/openjdk"
    },
    "data": [
        {
            "image": "docker.io/openjdk",
            "digest": "sha256:fe05457a5e9b9403f8e72eeba507ae80a4237d2d2d3f219fa62ceb128482a9ee",
            "processed": "2023-05-14T10:04:51.299231Z"
        }
    ]
}
```

## image version exposures

Get exposures for image version:

```shell
curl --data-urlencode "img=docker.io/openjdk" \
     --data-urlencode "dig=sha256:fe05457a5e9b9403f8e72eeba507ae80a4237d2d2d3f219fa62ceb128482a9ee" \
     --get https://vul.thingz.io/api/v1/exposures
```

Returns: 

```json
{
    "version": "v0.4.1",
    "created": "2023-05-14T12:45:07.36239582Z",
    "criteria": {
        "digest": "sha256:fe05457a5e9b9403f8e72eeba507ae80a4237d2d2d3f219fa62ceb128482a9ee",
        "image": "docker.io/openjdk"
    },
    "data": {
        "image": "docker.io/openjdk",
        "digest": "sha256:fe05457a5e9b9403f8e72eeba507ae80a4237d2d2d3f219fa62ceb128482a9ee",
        "packages": {
            "curl": {
                "versions": {
                    "7.61.1-25.el8_7.1": {
                        "sources": {
                            "grype": {
                                "exposures": {
                                    "CVE-2023-23916": {
                                        "severity": "medium",
                                        "score": 6.5,
                                        "fixed": true
                                    }
                                }
                            },
                            "snyk": {
                                "exposures": {
                                    "CVE-2023-23916": {
                                        "severity": "medium",
                                        "score": 6.5,
                                        "fixed": true
                                    }
                                }
                            },
                            "trivy": {
                                "exposures": {
                                    "CVE-2023-23916": {
                                        "severity": "medium",
                                        "score": 0,
                                        "fixed": true
                                    }
                                }
                            }
                        }
                    }
                }
            },
            ...
        }
    }
}
```


## Disclaimer

This is my personal project and it does not represent my employer. While I do my best to ensure that everything works, I take no responsibility for issues caused by this code.
