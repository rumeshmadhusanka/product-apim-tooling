/*
*  Copyright (c) WSO2 Inc. (http://www.wso2.org) All Rights Reserved.
*
*  WSO2 Inc. licenses this file to you under the Apache License,
*  Version 2.0 (the "License"); you may not use this file except
*  in compliance with the License.
*  You may obtain a copy of the License at
*
*    http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing,
* software distributed under the License is distributed on an
* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
* KIND, either express or implied.  See the License for the
* specific language governing permissions and limitations
* under the License.
 */

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wso2/product-apim-tooling/import-export-cli/operator/registry"
	k8sUtils "github.com/wso2/product-apim-tooling/import-export-cli/operator/utils"
	"github.com/wso2/product-apim-tooling/import-export-cli/utils"
)

const K8sInstallApiOperatorCmdLiteral = "api-operator"
const k8sInstallApiOperatorCmdShortDesc = "Install API Operator"
const k8sInstallApiOperatorCmdLongDesc = "Install API Operator in the configured K8s cluster"
const k8sInstallApiOperatorCmdExamples = utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral + `
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral + ` -f path/to/operator/configs
` + utils.ProjectName + ` ` + K8sCmdLiteral + ` ` + K8sInstallCmdLiteral + ` ` + K8sInstallApiOperatorCmdLiteral + ` -f path/to/operator/config/file.yaml`

// flags
var flagApiOperatorFile string

// flags for installing api-operator in batch mode
var flagBmRegistryType string
var flagBmRepository string
var flagBmUsername string
var flagBmPassword string
var flagBmPasswordStdin bool
var flagBmKeyFile string

// installApiOperatorCmd represents the install api-operator command
var installApiOperatorCmd = &cobra.Command{
	Use:     K8sInstallApiOperatorCmdLiteral,
	Short:   k8sInstallApiOperatorCmdShortDesc,
	Long:    k8sInstallApiOperatorCmdLongDesc,
	Example: k8sInstallApiOperatorCmdExamples,
	Run: func(cmd *cobra.Command, args []string) {
		utils.Logln(fmt.Sprintf("%s%s %s called", utils.LogPrefixInfo, K8sInstallCmdLiteral, K8sInstallApiOperatorCmdLiteral))

		// is -f or --from-file flag specified
		isLocalInstallation := flagApiOperatorFile != ""
		configFile := flagApiOperatorFile

		// check version before getting inputs (in interactive mode)
		if !isLocalInstallation {
			// getting API Operator version
			operatorVersion, err := k8sUtils.GetVersion(
				"API Operator",
				k8sUtils.ApiOperatorVersionEnvVariable,
				k8sUtils.DefaultApiOperatorVersion,
				k8sUtils.ApiOperatorVersionValidationUrlTemplate,
				k8sUtils.ApiOperatorFindVersionUrl,
			)
			if err != nil {
				utils.HandleErrorAndExit("Error in API Operator version", err)
			}
			configFile = fmt.Sprintf(k8sUtils.ApiOperatorConfigsUrlTemplate, operatorVersion)
		}

		// check for installation mode: interactive or batch mode
		// and get inputs
		if flagBmRegistryType == "" {
			// run api-operator installation in interactive mode
			// read inputs for docker registry
			registry.ChooseRegistryInteractive()
			registry.ReadInputsInteractive()
		} else {
			// run api-operator installation in batch mode
			// set registry type
			registry.SetRegistry(flagBmRegistryType)

			flagsValues := getGivenFlagsValues()
			registry.ValidateFlags(flagsValues)       // validate flags with respect to registry type
			registry.ReadInputsFromFlags(flagsValues) // read values from flags with respect to registry type
		}

		// installing operator and configs if -f flag given
		// otherwise settings configs only
		k8sUtils.CreateControllerConfigs(configFile, 20, k8sUtils.ApiOpCrdSecurity)
		registry.UpdateConfigsSecrets()

		fmt.Println("[Setting to K8s Mode]")
		utils.SetToK8sMode()
	},
}

// getGivenFlagsValues returns flags that user given in the batch mode except the "registry type"
func getGivenFlagsValues() *map[string]registry.FlagValue {
	flags := make(map[string]registry.FlagValue)
	flags[k8sUtils.FlagBmRepository] = registry.FlagValue{Value: flagBmRepository, IsProvided: flagBmRepository != ""}
	flags[k8sUtils.FlagBmUsername] = registry.FlagValue{Value: flagBmUsername, IsProvided: flagBmUsername != ""}
	flags[k8sUtils.FlagBmPassword] = registry.FlagValue{Value: flagBmPassword, IsProvided: flagBmPassword != ""}
	flags[k8sUtils.FlagBmPasswordStdin] = registry.FlagValue{Value: flagBmPasswordStdin, IsProvided: flagBmPasswordStdin}
	flags[k8sUtils.FlagBmKeyFile] = registry.FlagValue{Value: flagBmKeyFile, IsProvided: flagBmKeyFile != ""}

	return &flags
}

func init() {
	installCmd.AddCommand(installApiOperatorCmd)
	installApiOperatorCmd.Flags().StringVarP(&flagApiOperatorFile, "from-file", "f", "", "Path to API Operator directory")

	// flags for installing api-operator in batch mode
	// only the flag "registry-type" is required and others are registry specific flags
	installApiOperatorCmd.Flags().StringVarP(&flagBmRegistryType, "registry-type", "R", "", "Registry type: DOCKER_HUB | AMAZON_ECR |GCR | HTTP")
	installApiOperatorCmd.Flags().StringVarP(&flagBmRepository, k8sUtils.FlagBmRepository, "r", "", "Repository name or URI")
	installApiOperatorCmd.Flags().StringVarP(&flagBmUsername, k8sUtils.FlagBmUsername, "u", "", "Username of the repository")
	installApiOperatorCmd.Flags().StringVarP(&flagBmPassword, k8sUtils.FlagBmPassword, "p", "", "Password of the given user")
	installApiOperatorCmd.Flags().BoolVar(&flagBmPasswordStdin, k8sUtils.FlagBmPasswordStdin, false, "Prompt for password of the given user in the stdin")
	installApiOperatorCmd.Flags().StringVarP(&flagBmKeyFile, k8sUtils.FlagBmKeyFile, "c", "", "Credentials file")
}
