# Changelog

All notable changes to this project will be documented in this file. See [commit-and-tag-version](https://github.com/absolute-version/commit-and-tag-version) for commit guidelines.

## [1.0.1](https://github.com/dptsi/go-storage/compare/gcs/v1.0.0...gcs/v1.0.1) (2024-05-29)

## 1.0.0 (2024-05-29)


### ⚠ BREAKING CHANGES

* upload file by using interface instead of concrete

### Features

* **gcs:** upload, stream, and delete ([9f5871f](https://github.com/dptsi/go-storage/commit/9f5871fcaf812f98b6fbc4dc0abe3a97ecf82652))
* **s3:** base64 upload/download ([0c38dc4](https://github.com/dptsi/go-storage/commit/0c38dc446c00e9bcfadebc9f2b96ef35839e13eb))
* **s3:** delete file ([d569ece](https://github.com/dptsi/go-storage/commit/d569ece13f606ee4ecc50746ada58b43cb251812))
* **s3:** get file ([0749e37](https://github.com/dptsi/go-storage/commit/0749e371f91eb25d07cdb3ecaf2979d1606612cd))
* **s3:** sanitize file name method ([2c1ae89](https://github.com/dptsi/go-storage/commit/2c1ae89bd0211fc330f1ff80072fc97e7f1c5ac0))
* **s3:** temporary public link ([fbf1af5](https://github.com/dptsi/go-storage/commit/fbf1af507c34a682b46f797dee6524368a5428d6))
* **s3:** upload file ([e8ee587](https://github.com/dptsi/go-storage/commit/e8ee587c5b41254aa689e86d087759c0e445035d))
* **s3:** upload, download, delete, and get metadata, and file url ([5af7dea](https://github.com/dptsi/go-storage/commit/5af7deafe0f367757562457d1d56b36f902b6826))
* **storageapi:** delete file ([c6afbfc](https://github.com/dptsi/go-storage/commit/c6afbfc4e85fceb3acf32849dd04b30d3ac1675b))
* **storageapi:** get file by id ([3e81733](https://github.com/dptsi/go-storage/commit/3e8173379988985e5d020173d810af2ab0071e05))
* **storageapi:** upload file ([7c5a515](https://github.com/dptsi/go-storage/commit/7c5a515f54ca29a6a00cc535c72b37a44dba9a80))
* **storageapi:** upload, delete, and get file ([ef56baf](https://github.com/dptsi/go-storage/commit/ef56baf1a41b6b3752afd7f1155ef0685c93e511))
* upload file by using interface instead of concrete ([f332ee1](https://github.com/dptsi/go-storage/commit/f332ee15bab8d224cbe494e14a300f2911f3bb59))


### Bug Fixes

* **its:** can't upload file with name contains non-alphanumeric characters ([6c7b2da](https://github.com/dptsi/go-storage/commit/6c7b2da80ba9c698f520f3f8645e9261c7934897))
* wrong method signature for s3 ([6954cae](https://github.com/dptsi/go-storage/commit/6954cae947bf7c511d6d11438421e6c860e36c7b))
* wrong regex for sanitizing file name ([9c7625d](https://github.com/dptsi/go-storage/commit/9c7625d330555b682669cf604f8b0f38a1fc8a92))
