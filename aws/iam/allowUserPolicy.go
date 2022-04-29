package iam

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/mlibrodo/rds-db-copy/aws/config"
	"io"
	"strings"
)

// For explanation of these template variables see
// https://docs.aws.amazon.com/AmazonRDS/latest/UserGuide/UsingWithRDS.IAMDBAuth.IAMPolicy.html
type AssignDBToUserPolicy struct {
	Region        string
	AccountID     string
	DbiResourceId string
	DBUserName    string

	AWSUser      string
	DBInstanceID string
}

func (t *AssignDBToUserPolicy) createPolicy(out io.Writer) error {
	return templates.ExecuteTemplate(out, "iam_policy_allow_user_instance.tmpl", t)
}

func (t *AssignDBToUserPolicy) CreatePolicy() (*string, error) {
	sb := &strings.Builder{}

	if err := t.createPolicy(sb); err != nil {
		return nil, err
	}

	policy := sb.String()
	return &policy, nil
}

func (t *AssignDBToUserPolicy) AttachPolicyToUser() error {
	svc := iam.NewFromConfig(*config.AWSConfig)

	var policy *string
	var err error
	if policy, err = t.CreatePolicy(); err != nil {
		return err
	}

	policyName := fmt.Sprintf("DB=%s_User=%s", t.DBInstanceID, t.AWSUser)
	input := &iam.PutUserPolicyInput{
		PolicyDocument: policy,
		PolicyName:     &policyName,
		UserName:       &t.AWSUser,
	}

	if _, err = svc.PutUserPolicy(context.TODO(), input); err != nil {
		return err
	}
	return nil
}
