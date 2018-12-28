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
package io.elastest.ems.test.e2e;

import static io.github.bonigarcia.BrowserType.CHROME;
import static java.lang.invoke.MethodHandles.lookup;
import static org.openqa.selenium.support.ui.ExpectedConditions.visibilityOfElementLocated;
import static org.slf4j.LoggerFactory.getLogger;

import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInfo;
import org.junit.jupiter.api.extension.ExtendWith;
import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.remote.RemoteWebDriver;
import org.openqa.selenium.support.ui.WebDriverWait;
import org.slf4j.Logger;

import io.elastest.ems.test.base.EmsBaseTest;
import io.github.bonigarcia.BrowserType;
import io.github.bonigarcia.DockerBrowser;
import io.github.bonigarcia.SeleniumExtension;

/**
 * Check that the EMS works properly together with a TJob.
 */
@Tag("e2e")
@DisplayName("E2E tests of EMS through TORM")
@ExtendWith(SeleniumExtension.class)
public class EmsTJobE2eTest extends EmsBaseTest {
    final Logger log = getLogger(lookup().lookupClass());
    String projectName = "EMSe2e";

    private static final Map<String, List<String>> tssMap;
    static {
        tssMap = new HashMap<String, List<String>>();
        tssMap.put("EMS", null);
    }

    void createProject(WebDriver driver) throws Exception {
        navigateToTorm(driver);
        if (!etProjectExists(driver, projectName)) {
            createNewETProject(driver, projectName);
        }
    }

    private String sutName = "EMSe2esut";

    void createProjectAndSut(WebDriver driver) throws Exception {
        navigateToTorm(driver);
        if (!etProjectExists(driver, projectName)) {
            createNewETProject(driver, projectName);
        }
        if (!etSutExistsIntoProject(driver, projectName, sutName)) {
            // Create SuT
            String sutDesc = "SuT for E2E test";
            String sutMainService = "nginx-service";
            String sutSpec = "version: '3'\n\nservices:\n\n  nginx-service:\n   image: nginx\n   entrypoint:\n    - /bin/bash\n    - \"-c\"\n    - \"dd if=/dev/random of=/usr/share/nginx/html/sparse bs=1024 count=1 seek=5242880000;nginx;sleep infinity\"\n   expose:\n     - \"80\"";
            String sutPort = "80";
            createNewSutDeployedByElastestWithCompose(driver, sutName, sutDesc, sutSpec, sutMainService, sutPort, null, false);
        }

    }

    @Test
    @DisplayName("EMS in a TJob")
    void testTJob(@DockerBrowser(type = CHROME) RemoteWebDriver localDriver,
            TestInfo testInfo) throws Exception {
        setupTestBrowser(testInfo, BrowserType.CHROME, localDriver);

        // Setting up the TJob used in the test
        this.createProjectAndSut(driver);
        navigateToETProject(driver, projectName);
        String tJobName = "EMS e2e tjob";
        if (!etTJobExistsIntoProject(driver, projectName, tJobName)) {
            String tJobTestResultPath = "";
            String tJobImage = "imdeasoftware/e2etjob";
            String commands = "cd /go;./tjob";
            createNewTJob(driver, tJobName, tJobTestResultPath, sutName,
                    tJobImage, false, commands, null, tssMap, null);
        }
        // Run the TJob
        runTJobFromProjectPage(driver, tJobName);

        // Wait for eus card
        /*WebDriverWait waitEus = new WebDriverWait(driver, 60);
        By eusCard = By.xpath("//md-card-title[contains(string(), 'EUS')]");
        waitEus.until(visibilityOfElementLocated(eusCard));*/

        // and check its result
        this.checkFinishTJobExec(driver, 220, "SUCCESS", false);
    }
}
