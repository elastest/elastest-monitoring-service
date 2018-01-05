/*
 * (C) Copyright 2017-2019 ElasTest (http://elastest.io/)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */
package io.elastest.demo.rest;

import static org.assertj.core.api.Assertions.assertThat;

import org.junit.jupiter.api.Test;

import org.springframework.web.client.RestTemplate;

public class EmsTestAppTest {

    @Test
    public void rootServiceTest() {

        String appHost = System.getenv("ET_EMS_LSBEATS_HOST");
        if (appHost == null) {
            appHost = "172.27.0.9";
        }
		String ems_api = "http://" + appHost + ":8888/health";
		String ems_api_health = ems_api +  ":8888/health";
			
        RestTemplate client = new RestTemplate();
		
		/* 1. Configure an external elasticsearch+kibana web site for
		 * demoing purposes */

		/* 1.1 Create a query */
		
		// /* Method 1 */
		// String subscriber_request =
		// 	"{ \"channel\": \"any\","
		// 	+ "\"ip\": \"elastest.software.imdea.org\","
		// 	+ "\"port\": 9202,"
		// 	+ "\"user\": \"elastic\","
		// 	+ "\"password\": \"changeme\" }";

		// /* Method 2 */
		// JSONObject obj = new JSONObject();
		// obj.put("channel","any");
		// obj.put("ip","elastest.software.imdea.org");
		// obj.put("port",new Integer(9202));
		// obj.put("user","elastic");
		// obj.put("password","changeme");
		// StringWriter out = new StringWriter();
		// obj.writeJSONString(out);
		// String subscriber_request = out.toString();

		/* Method 3 */
		String subscriber_request =
			"{ 'channel': 'any',"
			+ "'ip': 'elastest.software.imdea.org',"
			+ "'port': 9202,"
			+ "'user': 'elastic',"
			+ "'password': 'changeme' }";


		/* 1.2 consume the API */
		String ems_api_subscribe = ems_api + "/subscriber/elasticsearch";


		System.out.println("Requesting EMS to subscribe: \"" + subscriber_request + "\"");
		String apiResponse = client.postForObject(ems_api_subscribe,
												  subscriber_request,
												  String.class);
        System.out.println("EMS responds: \"" + apiResponse + "\"");

        int counter = 60;

		int expected_events = 100;
		String result = "0";
        while (expected_events > Integer.parseInt(result) && counter > 0) {
			
			System.out.println("Connecting to "+ ems_api_health +" ...");
            result = client.getForObject(ems_api_health, String.class);

			System.out.println("Received: \"" + result + "\"");

            result = result.split(",")[1];
            result = result.split(":")[1];
            result = result.split("}")[0];

			System.out.println("...which corresponds to " + result + "events");

            counter--;
            try {
                System.out.println("sleeping for 3s...");
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            System.out.println("counter: " + counter + ". trying it again...");

        }
        assertThat(result).isNotEqualTo("0");
    }

}
