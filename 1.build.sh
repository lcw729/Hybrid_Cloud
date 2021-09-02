#!/bin/bash

controller_name="hcp-policy-engine"
resource_name="hcpPolicyEngine"
target_namespace="hybrid"


resource_name_small_letter=`echo $resource_name | tr '[A-Z]' '[a-z]'`


cp -r template-controller ./$controller_name

find $controller_name -type f -exec sed -i 's/template-controller/'$controller_name'/g' {} \;
find $controller_name -type f -exec sed -i 's/templateresource/'$resource_name_small_letter'/g' {} \;
find $controller_name -type f -exec sed -i 's/TemplateResource/'$resource_name'/g' {} \;
find $controller_name -type f -exec sed -i 's/nsnsns/'$target_namespace'/g' {} \;

mv $controller_name/pkg/controller/templateresource $controller_name/pkg/controller/$resource_name_small_letter

