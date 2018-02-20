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

import static java.lang.invoke.MethodHandles.lookup;
import static java.util.concurrent.TimeUnit.SECONDS;
import static org.openqa.selenium.support.ui.ExpectedConditions.textToBePresentInElementLocated;
import static org.openqa.selenium.support.ui.ExpectedConditions.visibilityOfElementLocated;
import static org.slf4j.LoggerFactory.getLogger;

import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.openqa.selenium.By;
import org.openqa.selenium.Dimension;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.support.ui.WebDriverWait;
import org.slf4j.Logger;

import io.elastest.ems.test.base.EmsBaseTest;
import io.github.bonigarcia.SeleniumExtension;

/**
 * E2E EMS test.
 *
 * @author Boni Garcia (boni.garcia@urjc.es)
 * @since 0.1.1
 */
@Tag("e2e")
@DisplayName("E2E tests of EMS through TORM")
@ExtendWith(SeleniumExtension.class)
public class EmsTJobE2eTest extends EmsBaseTest {

    final Logger log = getLogger(lookup().lookupClass());

    @Test
    @DisplayName("EMS in a TJob")
    void testTJob(ChromeDriver driver) throws InterruptedException {
        this.driver = driver;

        log.info("Navigate to TORM and start new project");
        driver.manage().window().setSize(new Dimension(1024, 1024));
        driver.manage().timeouts().implicitlyWait(5, SECONDS);
        if (secureElastest) {
            driver.get(secureTorm);
        }
        else {
            driver.get(tormUrl);
        }
        createNewProject(driver, "my-test-project");

        log.info("Create new TJob using EMS");
        driver.findElement(By.xpath("//button[contains(string(), 'New TJob')]"))
                .click();
        driver.findElement(By.name("tJobName")).sendKeys("my-test-tjob");
        driver.findElement(By.name("tJobImageName"))
                .sendKeys("elastest/ci-docker-e2e");
        driver.findElement(By.name("resultsPath")).sendKeys(
                "/home/jenkins/elastest-monitoring-service/tjob-test/target/surefire-reports/TEST-io.elastest.ems.test.e2e.TJobEmsTest.xml");
        driver.findElement(By.className("mat-select-trigger")).click();
        driver.findElement(By.xpath("//md-option[contains(string(), 'None')]"))
                .click();
        driver.findElement(By.name("commands")).sendKeys(
                "git clone https://github.com/fgorostiaga/elastest-monitoring-service; cd elastest-monitoring-service/tjob-test; mvn test;");
        driver.findElement(By.xpath("//md-checkbox[contains(string(), 'EUS')]"))
                .click();
        driver.findElement(By.xpath("//button[contains(string(), 'SAVE')]"))
                .click();

        log.info("Run TJob and wait for EMS GUI");
        driver.findElement(By.xpath("//button[@title='Run TJob']")).click();
        By eusCard = By
                .xpath("//md-card-title[contains(string(), 'elastest-eus')]");
        WebDriverWait waitEus = new WebDriverWait(driver, 60);
        waitEus.until(visibilityOfElementLocated(eusCard));

        log.info("Wait for build sucess traces");
        WebDriverWait waitLogs = new WebDriverWait(driver, 180);
        waitLogs.until(textToBePresentInElementLocated(By.tagName("logs-view"),
                "BUILD SUCCESS"));
    }

}
