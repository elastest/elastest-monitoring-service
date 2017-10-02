node('docker') {
    stage "Container Prep"
        echo("The node is up")
        def mycontainer = docker.image('elastest/ci-docker-compose-siblings')
        mycontainer.pull()
        mycontainer.inside("-u jenkins -v /var/run/docker.sock:/var/run/docker.sock:rw -v ${WORKSPACE}:/home/jenkins/.m2") {
            git 'https://github.com/elastest/elastest-monitoring-service'

            stage "Tests"
                echo ("Starting tests")
                sh 'docker run -v $(pwd)/go_EMS:/go/go_EMS golang /bin/bash -c "cd go_EMS;go test"'
                
            stage "Publish code coverage"
                echo ("Publishing code coverage")
                def codecovArgs = '-K '
                if (env.GITHUB_PR_NUMBER != '') {
                  // This is a PR
                  codecovArgs += "-B ${env.GITHUB_PR_TARGET_BRANCH} " +
                      "-C ${env.GITHUB_PR_HEAD_SHA} " +
                      "-P ${env.GITHUB_PR_NUMBER} "
                } else {
                  // Not a PR
                  codecovArgs += "-B ${env.GIT_BRANCH} " +
                      "-C ${env.GIT_COMMIT} "
                }
                sh "echo args = ${codecovArgs}"
                sh 'docker run -v $(pwd)/go_EMS:/go/go_EMS golang /bin/bash -c "cd go_EMS; go test -race -coverprofile=coverage.txt -covermode=atomic; curl -s https://codecov.io/bash | bash -s - ${codecovArgs} -t ${COB_EMS_TOKEN} || echo \'Codecov did not collect coverage reports\'"'

            stage "Build images - Package"
                echo ("Building full version")
                sh 'docker build -t elastest/ems:0.1 .'
                def myfullimage = docker.image('elastest/ems:0.1');
                echo ("Building min version")
                sh 'docker build -f Dockerfile_min -t elastest/ems_min:0.1 .'
                def myminimage = docker.image('elastest/ems_min:0.1');

            stage "Run images"
                myfullimage.run()
                myminimage.run()

            stage "Publish"
                echo ("Publishing")
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'elastestci-dockerhub',
                    usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                    sh 'docker login -u "$USERNAME" -p "$PASSWORD"'
                    myfullimage.push()
					myminimage.run()
                }   
        }   
}
