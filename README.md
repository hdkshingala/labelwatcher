# Label Watcher Kubernetes Operator

A kubernetes operator which validates Pods which are being edited and checks if the label from the Pod is not edited that is also used in a NetworkPolicy as a selector.

## How to run

- Run below command to generate Docker image:
    ```
    docker build -t <DOCKER_USERNAME>/labelwatcher:1.0.0 .
    ```

- Push the generated image to Docker hub:
    ```
    docker push <DOCKER_USERNAME>/labelwatcher:1.0.0
    ```

- Replace the image name in `manifest/install/deploy.yaml`.

- Run below command to deploy the controller on the Kubernetes:
    ```
    kubectl apply -f manifest/install/
    ```
