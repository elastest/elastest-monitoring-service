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

import static org.assertj.core.api.Assertions.*;
import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.ComparableAssert;

import org.junit.jupiter.api.Test;

import org.springframework.web.client.RestTemplate;

public class EmsTestAppTest {

    @Test
    public void rootServiceTest() {

        String appHost = System.getenv("ET_EMS_LSBEATS_HOST");
        if (appHost == null) {
            appHost = "172.27.0.9";
        }

        try {
            String endpointjson = "{ \"channel\": \"any\", \"ip\": \"elastest.software.imdea.org\", \"port\": 9202, \"user\": \"elastic\", \"password\": \"changeme\" }";
            Process p=Runtime.getRuntime().exec(new String[]{"curl","-d",endpointjson,"-H","Content-Type: application/json","http://"+appHost+":8888/subscriber/elasticsearch"});
            p.waitFor();
        } catch (Exception e) {
            assertThat(0).isNotEqualTo(0);
            System.out.println("exception: "+e);
        }

        RestTemplate client = new RestTemplate();

        String result = "0";
		int processed_events = 0;

        int counter = 60;

		int expected_events = 10000;
        while (expected_events > processed_events && counter > 0) {
			String ems_api_url = "http://" + appHost + ":8888/health";
			System.out.println("Connecting to "+ems_api_url+"...");
            result = client.getForObject(ems_api_url, String.class);

			System.out.println("Received: \"" + result + "\"");

            result = result.split(",")[1];
            result = result.split(":")[1];
            result = result.split("}")[0];

			System.out.println("...which corresponds to " + result + "events");

			processed_events = Integer.parseInt(result);

            counter--;
            try {
                System.out.println("sleeping for 3s...");
                Thread.sleep(3000);
            } catch (InterruptedException e) {
                e.printStackTrace();
            }
            System.out.println("counter: " + counter + ". trying it again...");

        }
		/* assertThat(result).isNotEqualTo("0");  */
		assertThat(processed_events).isGreaterThanOrEqual(expected_events); 
    }
}
