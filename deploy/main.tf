terraform {
  cloud {
    organization = "haraiai"

    workspaces {
      name = "haraiai"
    }
  }
}

provider "google" {
  project = "haraiai"
  region  = "asia-northeast1"
}

# Cloud Storage
resource "google_storage_bucket" "bucket" {
  name     = "haraiai-bucket"
  location = "US"
}

# Secret Manager
data "google_secret_manager_secret_version" "line_bot_channel_secret" {
  secret = "lineBotChannelSecret"
}

data "google_secret_manager_secret_version" "line_bot_channel_access_token" {
  secret = "lineBotChannelAccessToken"
}

# Cloud Functions for LINE Bot
data "archive_file" "bot" {
  type        = "zip"
  source_dir  = "../func/bot"
  output_path = "tmp/bot.zip"
}

resource "google_storage_bucket_object" "bot" {
  name   = "func/bot.${data.archive_file.bot.output_md5}.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.bot.output_path
}

resource "google_cloudfunctions_function" "bot_webhook" {
  name        = "BotWebhookHandler"
  description = "Handle LINE Bot webhook"
  runtime     = "go120"

  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.bot.name

  trigger_http = true
  entry_point  = "HandleWebhook"

  available_memory_mb = 128
  timeout             = 10
  min_instances       = 0
  max_instances       = 3

  environment_variables = {
    "PHASE"             = "production"
    "PROJECT_ID"        = "haraiai"
    "FE_BASE_URL"       = "https://haraiai.com"
    "PACKAGE_BASE_PATH" = "/workspace/serverless_function_source_code"

    "CHANNEL_SECRET"       = data.google_secret_manager_secret_version.line_bot_channel_secret.secret_data
    "CHANNEL_ACCESS_TOKEN" = data.google_secret_manager_secret_version.line_bot_channel_access_token.secret_data

    "TZ" = "Asia/Tokyo"
  }
}

resource "google_cloudfunctions_function_iam_member" "webhook_invoker" {
  project        = google_cloudfunctions_function.bot_webhook.project
  region         = google_cloudfunctions_function.bot_webhook.region
  cloud_function = google_cloudfunctions_function.bot_webhook.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

output "function_bot_webhook_url" {
  value = google_cloudfunctions_function.bot_webhook.https_trigger_url
}

