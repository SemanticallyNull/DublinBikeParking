# Dublin Bike Parking

[![CircleCI](https://circleci.com/gh/SemanticallyNull/DublinBikeParking.svg?style=svg)](https://circleci.com/gh/SemanticallyNull/DublinBikeParking)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FSemanticallyNull%2FDublinBikeParking.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2FSemanticallyNull%2FDublinBikeParking?ref=badge_shield)

## Configuration

All the configuration is done using environment variables

| Environment Variable | Required | Default Value | Description |
| --- | --- | --- | --- |
| `DBP_DB_DIALECT` | false | `sqlite3` | The dialect to use. It can be one of `sqlite3` or `mysql` |
| `DBP_DB_CONNECTION_STRING` | false | `./demo.db` | The connection string for the DB. See [GORM's documentation](https://gorm.io/docs/connecting_to_the_database.html) for details. Only sqlite3 or mysql are currently supported. |
| `SENDGRID_API_KEY` | false | none | The project uses Sendgrid to send emails when a new stand is added for approval. If set to empty a warning will be printed, but the application will just not send email. |
| `DBP_UI_V2` | false | `false` | The new Vue UI is under development, this flag turns it on. |
| `S3_ENDPOINT` | false | none | S3 endpoint |
| `S3_ACCESS_KEY_ID` | false | none | S3 access key ID |
| `S3_SECRET_ACCESS_KEY` | false | none | S3 secret access key |
| `S3_BUCKET_NAME` | false | none | S3 bucket name |
| `OIDC_AUTHORITY` | false | none | The domain of the OIDC provider |
| `OIDC_AUDIENCE` | false | none | The audience for which the JWT Token is validated against |

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2FSemanticallyNull%2FDublinBikeParking.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2FSemanticallyNull%2FDublinBikeParking?ref=badge_large)