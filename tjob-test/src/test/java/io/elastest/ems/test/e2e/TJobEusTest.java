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
package io.elastest.eus.test.e2e;

import static java.lang.System.getenv;
import static java.lang.invoke.MethodHandles.lookup;
import static org.hamcrest.CoreMatchers.equalTo;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.openqa.selenium.remote.DesiredCapabilities.chrome;
import static org.slf4j.LoggerFactory.getLogger;

import java.net.MalformedURLException;
import java.net.URL;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.Capabilities;
import org.openqa.selenium.remote.RemoteWebDriver;
import org.slf4j.Logger;

/**
 * TJob test.
 *
 * @author Boni Garcia (boni.garcia@urjc.es)
 * @since 0.1.1
 */
public class TJobEusTest {

    final Logger log = getLogger(lookup().lookupClass());

    RemoteWebDriver driver;

    @BeforeEach
    void setup() throws MalformedURLException {
        Capabilities capabilities = chrome();
        String driverUrl = getenv("ET_EUS_API");
        if (driverUrl == null) {
            driverUrl = "http://172.21.0.10:8040/eus/v1/";
        }
        log.info("Using EUS URL {}", driverUrl);
        driver = new RemoteWebDriver(new URL(driverUrl), capabilities);
    }

    @Test
    void tJobTest() {
        driver.get("http://elastest.io/");
        assertThat(driver.getTitle(), equalTo("ElasTest Home"));
    }

    @AfterEach
    void teardown() {
        if (driver != null) {
            driver.quit();
        }
    }

}
