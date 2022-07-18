# Handy commands

```
rm -rf main && go build main.go && zip function.zip main
aws lambda update-function-code --function-name course-manager --zip-file fileb://function.zip --publish
aws lambda invoke --function-name course-manager --payload '{ "path": "/v1/listCourses" }' --cli-binary-format raw-in-base64-out output.json && cat output.json
```