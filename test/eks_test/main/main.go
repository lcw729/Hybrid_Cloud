package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
)

func main() {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
	}))
	eksSvc := eks.New(sess, aws.NewConfig().WithRegion("us-east-2"))
	input := &eks.CreateClusterInput{
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds: []*string{
				aws.String("subnet-0ea4f397a0207f679"),
				aws.String("subnet-0794e08bd9f4f1de4"),
				aws.String("subnet-0ee1e8932a11517f7"),
				aws.String("subnet-02fabab8ef21f3297"),
			},
		},
		Name:    aws.String("eks-cluster"),
		RoleArn: aws.String("arn:aws:iam::741566967679:role/eksClusterRole"),
	}

	result, err := eksSvc.CreateCluster(input)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

	result2, err := eksSvc.CreateNodegroup(&eks.CreateNodegroupInput{
		NodegroupName: aws.String("nodegroup1"),
		ClusterName:   aws.String("eks-cluster"),
		NodeRole:      aws.String("arn:aws:iam::741566967679:role/AmazonEKSNodeRole"),
		InstanceTypes: []*string{aws.String("t3.micro")},
		Subnets: []*string{
			aws.String("subnet-0ea4f397a0207f679"),
			aws.String("subnet-0794e08bd9f4f1de4"),
			aws.String("subnet-0ee1e8932a11517f7"),
			aws.String("subnet-02fabab8ef21f3297"),
		},
	})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result2)
	}
}
