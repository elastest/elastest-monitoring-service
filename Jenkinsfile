node('docker') {
    stage "Container Prep"
        echo("The node is up")
        def mycontainer = docker.image('elastest/ci-docker-compose-siblings')
        mycontainer.pull()
        mycontainer.inside("-u jenkins -v /var/run/docker.sock:/var/run/docker.sock:rw -v ${WORKSPACE}:/home/jenkins/.m2") {
            git 'https://github.com/elastest/elastest-monitoring-service'

            stage "Test and publish code coverage"
                echo ("Publishing code coverage")
                sh "mkdir shared || true"
                sh 'export PWD=$(pwd)'
                sh 'docker run -v ${PWD}/shared:/shared -v ${PWD}/striver-go:/go/src/gitlab.software.imdea.org/felipe.gorostiaga/striver-go -v ${PWD}:/go/src/github.com/elastest/elastest-monitoring-service golang /bin/bash -c "go get github.com/golang/protobuf/proto; go get google.golang.org/grpc; go get google.golang.org/grpc/reflection; cd src/github.com/elastest/elastest-monitoring-service/go_EMS; go test ./... -race -coverprofile=coverage.txt -covermode=atomic; mv coverage.txt /shared"'
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
