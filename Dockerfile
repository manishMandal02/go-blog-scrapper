# base image - stage 0
FROM golang:1.22-alpine AS build
# current working dir
WORKDIR /app                            
COPY . .                                    
# Install necessary dependencies and build 
RUN apk --no-cache add gcompat build-base \
    && wget https://github.com/tailwindlabs/tailwindcss/releases/latest/download/tailwindcss-linux-x64 \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 tailwindcss \
    && go install github.com/a-h/templ/cmd/templ@latest \
    && make build  


# runtime image
FROM alpine:3.18 AS final
# isntall chromium
RUN apk update && apk upgrade \
    && apk add --no-cache chromium
# workdir in the runtime image
WORKDIR /app                            
# copy contents from stage 0
COPY --from=build /app ./      

# Set environment variables
ENV CHROME_PATH=/usr/bin/chromium-browser

# Ensure the binary has executable permissions
RUN chmod +x ./bin/scrapper

# run the server
CMD ["./bin/scrapper"]                   


