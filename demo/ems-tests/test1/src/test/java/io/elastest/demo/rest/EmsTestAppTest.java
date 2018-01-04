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

        RestTemplate client = new RestTemplate();

        String result = "0";

        int counter = 60;

        while ("0".equals(result) && counter > 0) {
            result = client.getForObject("http://" + appHost + ":8888/health",
                    String.class);

            result = result.split(",")[1];
            result = result.split(":")[1];
            result = result.split("}")[0];

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
