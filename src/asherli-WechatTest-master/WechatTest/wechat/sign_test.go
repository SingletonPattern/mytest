/*
 *    Copyright 2016-2018 Li ZongZe
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package wechat

import "testing"

func TestValidateURL(t *testing.T) {
	token := "0t37dWsIYg6NsVLgEY1fNuB1rSLyyeQE"
	timestamp := "1449648662"
	nonce := "1862651475"
	signature := "717efa7b4910821c7bd59c1b84bbfc363f7551ef"

	ok := ValidateURL(token, timestamp, nonce, signature)

	if !ok {
		t.Fail()
	}
}

func TestSignature(t *testing.T) {
	token := "0t37dWsIYg6NsVLgEY1fNuB1rSLyyeQE"
	timestamp := "1449648662"
	nonce := "18626514725"
	encrypt := "lvMjItfR0rOPRpWGTG3K/b6zEKg4HDKeMU+/HtH6xqZJPpO0fQS8aSVmtornTIowI394/0xSjfxUNT7fdEJvGYpbgU0c2S8P8fQ/+oinc73tEl1hCJSsButo8tPYhjzKzuVITf9OSw4AcS7oo8W8SQBW5ndhj/Cy//kkRm4B82luwpTGHJ8RVcwXriGHVcnW56tYNnmgbGDie2Y0o3vkXX2Gvl7x0iDQpl8QgenMDm4OhvmAL5irMUtPiCFqvB1YM9LCN/f5dbwxMXFdjcI1XJIc6pY6e3t5SC9v96bH+UxgGls5IQuA/ZjNQOFREUp6G3S9A2cvRiNd/jjI72kLbl10KcJRotw1ozkDL8q+azT1OqisNQecsrC/sJ915FlNXbRzSI14RA9HWEOi8XCvphmJlLcoSIYQ/YyC70724tg="
	msgSign := "1f06de4929383d3244fe6981b75dc12ca4928aa2"
	actual := Signature(token, timestamp, nonce, encrypt)

	if actual != msgSign {
		t.Fail()
	}
}
