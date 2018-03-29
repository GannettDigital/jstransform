#!groovy

node {
  stage 'Execute Docker CI'

  // paas-api-ci options
  def paasApiCiVersion = "5.8.14-156"
  def repo = "GannettDigital/jstransform"
  def environment = "staging"
  def region = "us-east-1"
  def vaultURL = "https://vault.service.us-east-1.gciconsul.com:8200"
  def vaultConfig = "/secret/paas-api/paas-api-ci"

  withCredentials([
    string(credentialsId: "X_API_KEY", variable: "X_API_KEY"),
    string(credentialsId: "X_TYK_API_KEY", variable: "X_TYK_API_KEY"),
    string(credentialsId: "SLACK_URL", variable: "SLACK_URL"),
    string(credentialsId: "PAAS_API_CI_VAULT_TOKEN", variable: "PAAS_API_CI_VAULT_TOKEN")
  ]) {
    try {

      print 'Running docker run'

      sh "docker run -e \"GIT_BRANCH=${env.BRANCH_NAME}\" -e \"BUILD_ID=${env.BUILD_ID}\" -e \"VAULT_ADDR=${vaultURL}\" -e \"VAULT_CONFIG_LOCATION=${vaultConfig}\" -e \"VAULT_TOKEN=${PAAS_API_CI_VAULT_TOKEN}\" -e X_TYK_API_KEY=\"${X_TYK_API_KEY}\" -e X_API_KEY=\"${X_API_KEY}\" --rm -v /var/run/docker.sock:/var/run/docker.sock -v ~/.docker/config.json:/root/.docker/config.json paas-docker-artifactory.gannettdigital.com/paas-api-ci:${paasApiCiVersion} build \
        --repo=\"${repo}\" \
        --package \
        --skip-deploy"
    } catch (err) {
      currentBuild.result = "FAILURE"

      if (env.JOB_NAME.contains("master")) {
        def slackNotify = 'curl -X POST --data-urlencode \'payload= { "channel": "#api-releases", "username": "gopher-bot", "icon_emoji": ":httpmock:", "text": "*Build Failure*", "attachments": [ { "color": "#ff0000", "text": "<!here> Build failure for '
        slackNotify += env.JOB_NAME
        slackNotify += '","fields": [{"title": "Build URL","value": "'
        slackNotify += env.BUILD_URL
        slackNotify += '","short": true},{"title": "Build Number","value": "'
        slackNotify += env.BUILD_NUMBER
        slackNotify += '","short": true}]}]}\' '
        slackNotify += SLACK_URL

        sh slackNotify
      }

      throw err
    }
  }
}
