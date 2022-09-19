package main

import (
	"github.com/blushft/go-diagrams/diagram"
	"github.com/blushft/go-diagrams/nodes/aws"
	"log"
)

func main() {
	d, err := diagram.New(diagram.Filename("diagram"), diagram.Label("aws-playground-go"))
	if err != nil {
		log.Fatal(err)
	}

	ec2 := aws.Compute.Ec2().Label("Pet Store API")
	db := aws.Database.Rds().Label("Database")
	s3 := aws.Storage.SimpleStorageServiceS3().Label("S3 bucket")
	lambda := aws.Compute.Lambda().Label("Generate image thumbnail")
	cw := aws.Management.Cloudwatch().Label("Cloudwatch Alarm")

	dc := diagram.NewGroup("AWS")

	d.Connect(ec2, s3, diagram.Forward(), label("Upload pet image")).Group(dc)
	d.Connect(s3, cw, diagram.Forward(), label("On PUT trigger")).Group(dc)
	d.Connect(cw, lambda, diagram.Forward(), label("Trigger Lambda")).Group(dc)
	d.Connect(lambda, s3, diagram.Forward(), label("Upload thumbnail")).Group(dc)
	d.Connect(ec2, db, diagram.Forward(), label("Save pet info")).Group(dc)

	if err := d.Render(); err != nil {
		log.Fatal(err)
	}

	// run: dot -Tpng diagram.dot > diagram.png
}

func label(label string) func(o *diagram.EdgeOptions) {
	return func(o *diagram.EdgeOptions) {
		o.Label = label
	}
}
