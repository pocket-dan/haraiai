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

data "google_secret_manager_secret_version" "line_notifier_receiver_token" {
  secret = "lineNotifierReceiverToken"
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
  runtime     = "go116"

  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.bot.name

  trigger_http = true
  entry_point  = "HandleWebhook"

  available_memory_mb = 128
  timeout             = 10
  min_instances       = 0
  max_instances       = 3

  environment_variables = {
    "PHASE" = "production"
    "PROJECT_ID" = "haraiai"
    "FE_BASE_URL" = "https://haraiai.netlify.app"
    "PACKAGE_BASE_PATH" = "/workspace/serverless_function_source_code"

    "CHANNEL_SECRET" = data.google_secret_manager_secret_version.line_bot_channel_secret.secret_data
    "CHANNEL_ACCESS_TOKEN" = data.google_secret_manager_secret_version.line_bot_channel_access_token.secret_data
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


# Cloud Functions for general REST API
data "archive_file" "api" {
  type        = "zip"
  source_dir  = "../func/api"
  output_path = "tmp/api.zip"
}

resource "google_storage_bucket_object" "api" {
  name   = "func/api.${data.archive_file.api.output_md5}.zip"
  bucket = google_storage_bucket.bucket.name
  source = data.archive_file.api.output_path
}

resource "google_cloudfunctions_function" "api_notify_inquiry" {
  name        = "NotifyInquiry"
  description = "Handle inquiry from user"
  runtime     = "go116"

  source_archive_bucket = google_storage_bucket.bucket.name
  source_archive_object = google_storage_bucket_object.api.name

  trigger_http = true
  entry_point  = "NotifyInquiry"

  available_memory_mb = 128
  timeout             = 10
  min_instances       = 0
  max_instances       = 1

  environment_variables = {
    "PHASE" = "production"
    "FE_BASE_URL" = "https://haraiai.netlify.app"

    "LINE_NOTIFY_TOKEN" = data.google_secret_manager_secret_version.line_notifier_receiver_token.secret_data
  }
}

resource "google_cloudfunctions_function_iam_member" "inquiry_invoker" {
  project        = google_cloudfunctions_function.api_notify_inquiry.project
  region         = google_cloudfunctions_function.api_notify_inquiry.region
  cloud_function = google_cloudfunctions_function.api_notify_inquiry.name

  role   = "roles/cloudfunctions.invoker"
  member = "allUsers"
}

output "function_api_inquiry_url" {
  value = google_cloudfunctions_function.api_notify_inquiry.https_trigger_url
}

