node('docker') {
    stage "Container Prep"
        echo("The node is up")
        def mycontainer = docker.image('elastest/ci-docker-compose-siblings')
        mycontainer.pull()
        mycontainer.inside("-u jenkins -v /var/run/docker.sock:/var/run/docker.sock:rw -v ${WORKSPACE}:/home/jenkins/.m2") {
            git 'https://github.com/elastest/elastest-monitoring-service'

            stage "Tests"
                echo ("Starting tests")
                sh 'docker run -v $(pwd)/go_EMS:/go/src/github.com/elastest/elastest-monitoring-service/go_EMS golang /bin/bash -c "cd src/github.com/elastest/elastest-monitoring-service/go_EMS;go test"'
                
            stage "Publish code coverage"
                echo ("Publishing code coverage")
                sh "mkdir shared || true"
                sh 'export PWD=$(pwd)'
                sh 'docker run -v ${PWD}/shared:/shared -v ${PWD}/go_EMS:/go/go_EMS golang /bin/bash -c "cd go_EMS; go test -race -coverprofile=coverage.txt -covermode=atomic; mv coverage.txt /shared"'
                sh "curl -s https://codecov.io/bash > shared/curlout.txt"
                sh "cd shared; JENKINS_URL= bash <curlout.txt -s - -t ${COB_EMS_TOKEN}; cd ..; rm -rf shared"

            stage "Build images - Package"
                echo ("Building full version")
                sh 'docker build --build-arg GIT_COMMIT=$(git rev-parse HEAD) --build-arg COMMIT_DATE=$(git log -1 --format=%cd --date=format:%Y-%m-%dT%H:%M:%S) . -t elastest/ems:latest'
                def myfullimage = docker.image('elastest/ems:latest');

            stage "Run images"
                myfullimage.run()

            stage "Publish"
                echo ("Publishing")
                withCredentials([[$class: 'UsernamePasswordMultiBinding', credentialsId: 'elastestci-dockerhub',
                    usernameVariable: 'USERNAME', passwordVariable: 'PASSWORD']]) {
                    sh 'docker login -u "$USERNAME" -p "$PASSWORD"'
                    myfullimage.push()
                }
        }   
}
