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

package apim

// ApplicationList : Application List DTO
type ApplicationList struct {
	Count int               `json:"count"`
	List  []ApplicationInfo `json:"list"`
}

// ApplicationInfo : Application info DTO
type ApplicationInfo struct {
	ApplicationID     string   `json:"applicationId"`
	Name              string   `json:"name"`
	ThrottlingPolicy  string   `json:"throttlingPolicy"`
	Description       string   `json:"description"`
	Status            string   `json:"status"`
	Groups            []string `json:"groups"`
	SubscriptionCount int      `json:"subscriptionCount"`
	Owner             string   `json:"owner"`
}
