echo "Removing Docker image, Building Docker image..." && docker image rm ioxposter && docker buildx build --platform linux/arm/v7 -t ioxposter . 
