set -eo pipefail
cd function
go test
GOOS=linux go build main.go
cd ../
aws cloudformation package --template-file template.yml --s3-bucket $ARTIFACT_BUCKET --output-template-file out.yml
aws cloudformation deploy --template-file out.yml --stack-name pet-store --capabilities CAPABILITY_NAMED_IAM  --region $AWS_REGION