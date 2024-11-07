CURRENTDATE=`date +"%Y%m%d%H%M%S"`

echo "Generating Swagger Files"
$HOME/go/bin/swag init

echo "Building Executable"
go mod tidy
GOOS=linux GOARCH=amd64 go build .

echo "Building Docker Image"
docker build --platform linux/amd64 --tag jacohanekom/sanitize:latest .
rm sanitize

echo "Saving Image"
rm -rf work
mkdir -p work
docker save jacohanekom/sanitize:latest > work/sanitize.tar

docker image rm jacohanekom/sanitize:latest

echo "Packing Build"

cp docker-compose.yml work/docker-compose.yml
cp docs/swagger.json work/swagger.json
cp docs/swagger.yaml work/swagger.yaml
mkdir -p releases
tar -zcvf releases/${CURRENTDATE}.tar.gz -C work/ .
rm -rf work

