terraform {
  required_providers {
    archive = {
      source  = "hashicorp/archive"
      version = "2.2.0"
    }

    google = {
      source  = "hashicorp/google"
      version = "5.27.0"
    }
  }

  required_version = ">= 1.1"
}
