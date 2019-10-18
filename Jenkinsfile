/* FIXME: This uses Docker build commands via shell instead of the native
   Docker Pipeline plugin DSL to work around multi-stage build failures per
   this issue: https://issues.jenkins-ci.org/browse/JENKINS-44609 .
   I chose to invoke unit tests from within the Dockerfile to simplify
   local development and testing in a build-like environment, while avoiding
   duplicate test definitions in Jenkins. There's got to be a better way! */

pipeline {
    agent any
    options {
        skipStagesAfterUnstable()
    }
    environment {
      // This allows us to compose a build revision from the shortened branch name, short commit hash, and build number.
      BUILD = "${env.BRANCH_NAME}-${env.GIT_COMMIT[0..6]}-${env.BUILD_NUMBER}"
      JOB_ROOT = "aphorismophilia" 
    }
    stages {
        stage('Unlock Repository'){
            environment {
                GIT_CRYPT_KEY = credentials('GitCrypt_Key_aphorismophilia')
            }
            steps {
                sh label: 'Reset git-crypt status', script: 'git-crypt lock --force', returnStatus: true
                sh label: 'Decrypt repository with git-crypt', script: 'git-crypt unlock ${GIT_CRYPT_KEY}'
            }
        }
        stage('Go Test') {
            steps {
                sendNotifications 'STARTED'
                sh label: 'Build test container and run Go unit tests (inline Dockerfile RUN cmds)', script: 'docker build --target test --build-arg BUILD=${BUILD} -t ${JOB_ROOT}-test:${BUILD} .'
                sh label: 'Retrieve test results from test container', script: 'docker run --rm --entrypoint cat ${JOB_ROOT}-test:${BUILD} /go/src/aphorismophilia/ut-results.xml > ut-results-${BUILD}.xml'
            }
            post {
                always {
                    junit testResults: "ut-results-${env.BUILD}.xml"
                }
                cleanup {
                    sh label: 'Clean up test results from workspace', script: 'rm -v ut-results-${BUILD}.xml'  
                }
            }
        }
        stage('Build') {
            steps {
                // NB: This reruns the unit tests too. Perhaps reorder Dockerfile, or tag+push the test/build images to registry and pull back in here. Revisit after JENKINS-44609 is fixed.
                sh label: 'Build release container', script: 'docker build --target release --build-arg BUILD=${BUILD} -t ${JOB_ROOT}:${BUILD} .'
            }
        }
        stage('Blackbox HTTP Test') {
            steps {
                sh label: 'Launch release container for blackbox testing', script: 'docker run -d --rm --name ${JOB_ROOT}-${BUILD}-blackbox ${JOB_ROOT}:${BUILD}'
                sh label: 'Run HTTP check against base service', script: '/usr/lib/nagios/plugins/check_http -H `docker inspect --format "{{ .NetworkSettings.IPAddress }}" ${JOB_ROOT}-${BUILD}-blackbox` -p 8888 -v -u "/"'
                sh label: 'Run HTTP check against fortune backend', script: '/usr/lib/nagios/plugins/check_http -H `docker inspect --format "{{ .NetworkSettings.IPAddress }}" ${JOB_ROOT}-${BUILD}-blackbox` -p 8888 -v -u "/?backend=fortune"'
                sh label: 'Run HTTP check against flatfile backend', script: '/usr/lib/nagios/plugins/check_http -H `docker inspect --format "{{ .NetworkSettings.IPAddress }}" ${JOB_ROOT}-${BUILD}-blackbox` -p 8888 -v -u "/?backend=flatfile"'
            }
            post {
                always {
                    sh label: 'Kill release container after blackbox testing', script: 'docker kill ${JOB_ROOT}-${BUILD}-blackbox'
                }
            }
        }
        /* Launch a prod-like integration testing namespace into the auto environment for pull request
           and master branch builds. This could launch a fully prod-consistent composed infrastructure
           stack, but I don't want to spend that time and money on personal project pipeline builds
           (at least while this application's dependencies still fit inside a Kubernetes namespace)
           especially since I already launch new stacks via the infrastructure template pipeline. */
        stage('Integration Test') {
            when {
                anyOf { changeRequest() ; branch 'master'}
            }
            environment { // These variables are passed via make to retrieve K8s provider credentials and invoke Terraform.
                CLOUDSDK_AUTH_CREDENTIAL_FILE_OVERRIDE = "../secrets/ephemeral-testing-service-account.json"
                KUBECONFIG = "../secrets/ephemeral-kubeconfig"
                A13A_VERSION = "${env.BUILD}-integration"
            }
            steps {
                dir("./terraform") {
                    withDockerRegistry([url: "", credentialsId: "DockerHub_mikeroach"]) {
                        sh label: 'Tag integration test container with registry repository info', script: 'docker tag ${JOB_ROOT}:${BUILD} mikeroach/${JOB_ROOT}:${BUILD}-integration'
                        sh label: 'Push integration test container to registry', script: 'docker push mikeroach/${JOB_ROOT}:${BUILD}-integration'
                    }
                    sh label: 'Retrieve K8s cluster credentials', script: 'make k8s-credentials'
                    sh label: 'Launch ephemeral integration testing environment', script: 'make environment'
                    sleep(time:15, unit:"SECONDS") // Wait 15 seconds for integration test pod(s) to launch.
                    sh label: 'Run HTTP check against integration environment', script: '/usr/lib/nagios/plugins/check_http -H `make -s http-host` -v -u "/" -s "Build: ${BUILD}"'
                  }
            }
            post {
                cleanup {
                    dir("./terraform") {
                        sh label: 'Destroy ephemeral integration testing environment', script: 'make destroy', returnStatus: true
                        sh label: 'Destroy ephemeral K8s cluster credentials', script: 'rm -f ${KUBECONFIG}', returnStatus: true
                    }
                }
            }
        }
        stage('Publish to Container Registry') {
            when {
                branch 'master' // Only push to container registry when updating master
            }
            steps {
                withDockerRegistry([url: "", credentialsId: "DockerHub_mikeroach"]) {
                    sh label: 'Tag release container with registry repository info', script: 'docker tag ${JOB_ROOT}:${BUILD} mikeroach/${JOB_ROOT}:${BUILD}'
                    sh label: 'Push release container to registry', script: 'docker push mikeroach/${JOB_ROOT}:${BUILD}'
                }
            }
        }
        stage('Update Auto Environments') {
          when {
              branch 'master' // Only deploy to upstream environment when updating master
          }
          steps { // I'll keep my local minikube test commands here for posterity.
              // Note the escaped \ in sed regexp match group; \ is a Groovy DSL special character 
              //sh label: 'Update live Kubernetes deployment', script: '~jenkins/minikube/kubectl --kubeconfig ~jenkins/minikube/config get deploy/${JOB_ROOT} -o yaml | sed "s/\\(image: registry.hub.docker.com\\/mikeroach\\/${JOB_ROOT}\\):.*$/\\1:${BUILD}/" | ~jenkins/minikube/kubectl --kubeconfig ~jenkins/minikube/config apply -f -'
              //sh label: 'Verify Kubernetes deployment completes', script: '~jenkins/minikube/kubectl --kubeconfig ~jenkins/minikube/config rollout status -w deploy/${JOB_ROOT}'
              //sh label: 'Verify new build deployed to live environment', script: '~jenkins/bin/check_http -H k8shost.local -v -u "/" -s "Build: ${BUILD}"'
              withCredentials([usernamePassword(credentialsId: 'GitHub_Jenkins-GCP', passwordVariable: 'GIT_PASSWORD', usernameVariable: 'GIT_USERNAME')]) {
                  sh label: 'Update auto pipeline with requested version', script: '''
                      git config --local user.email "jenkins@borrowingcarbon.net" && git config --local user.name "Jenkins"
                      git config --local credential.helper "!p() { echo username="$GIT_USERNAME" ; echo password="$GIT_PASSWORD" ; }; p"
                      rm -rf auto-pipeline/
                      git clone "https://${GIT_USERNAME}:${GIT_PASSWORD}@github.com/mikeroach/iac-pipeline-auto.git" ./auto-pipeline
                      cd auto-pipeline
                      GITHUB_TOKEN=${GIT_PASSWORD} ./gitops-helper.sh -p -s service_aphorismophilia -v ${BUILD} -u ${BUILD_URL}
                  '''
                }
            }
        }
    }
    post {
      success {
            script {
                if (env.CHANGE_ID) { // Only interact with pull requests if this job was triggered by one
                    /* Because we merge and close the pull request with the GitHub Pipeline plugin during this
                       job's post-success script, the Branch Source plugin can't update the pr-merge status
                       check on GitHub. To solve this we'll set test status from within the post script. */
                    pullRequest.createStatus(status: 'success',
                     context: 'continuous-integration/jenkins/pr-merge',
                     description: 'All tests passed',
                     targetUrl: "${env.BUILD_URL}/display/redirect")

                    // Attempt to auto-merge this PR into master unless the 'no-merge' label exists to indicate otherwise.
                    if (! pullRequest.labels.contains("no-merge")) {
                        echo "No-merge label absent; attempting auto-merge."

                        // Merge and close the PR if conflict-free, otherwise leave it open with a comment.
                        if (pullRequest.mergeable) {
                            pullRequest.comment('[Jenkins] All tests for this PR succeeded, merging and closing.')
                            pullRequest.merge('[Jenkins] Automatically merged by Jenkins.')
                        } else {
                            pullRequest.comment('[Jenkins] All tests for this PR succeeded, but PR is unmergeable. Please investigate.')
                        }
                    } else {
                        echo "No-merge label detected; skipping auto-merge."
                        pullRequest.comment('[Jenkins] All tests for this PR succeeded, skipping merge due to "no-merge" label.')
                    }
                }
            }
        }
        cleanup {
            sh label: 'Lock repository secret files', script: 'git-crypt lock --force', returnStatus: true
            sendNotifications currentBuild.result
        }
    }
}

// Retrieve changelog for notifications adapted from https://support.cloudbees.com/hc/en-us/articles/217630098-How-to-access-Changelogs-in-a-Pipeline-Job-
def getChangeString() {
 MAX_MSG_LEN = 100
 def changeString = ""

 echo "Gathering SCM changes"
 def changeLogSets = currentBuild.changeSets
 for (int i = 0; i < changeLogSets.size(); i++) {
 def entries = changeLogSets[i].items
 for (int j = 0; j < entries.length; j++) {
 def entry = entries[j]
 truncated_msg = entry.msg.take(MAX_MSG_LEN)
 changeString += " - ${truncated_msg} [${entry.author}]\n"
 }
 }

 if (!changeString) {
 changeString = " - No new changes"
 }
 return changeString
}

// Send Slack notifications adapted from https://jenkins.io/blog/2017/02/15/declarative-notifications/
/* MIT License

Copyright (c) 2017 Liam Newman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE. */
def sendNotifications(String buildStatus = 'STARTED') {
  // build status of null means successful
  buildStatus = buildStatus ?: 'SUCCESS'

  // Default values
  def colorName = 'RED'
  def colorCode = '#FF0000'
  def emoji = ":question:"
  def subject = "${buildStatus}: Job '${env.JOB_NAME} [${env.BUILD}]'"
  def summary = "${subject} (${env.BUILD_URL})"
    
  // Override default values based on build status
  if (buildStatus == 'STARTED') {
    emoji = ":construction:"
    color = 'YELLOW'
    colorCode = '#FFFF00'
    summary = "${emoji} ${subject} (${env.BUILD_URL})\n " + getChangeString()
  } else if (buildStatus == 'SUCCESS') {
    emoji = ":white_check_mark:"
    color = 'GREEN'
    colorCode = '#00FF00'
    summary = "${emoji} ${subject} (${env.BUILD_URL})"
  } else {
    emoji = ":x:"
    color = 'RED'
    colorCode = '#FF0000'
    summary = "${emoji} ${subject} (${env.BUILD_URL})" 
  }

  // Send notifications
  slackSend (channel: "roachtest", color: colorCode, message: summary)
}