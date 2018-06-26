/*
 * (C) Copyright 2017 OpenVidu (http://openvidu.io/)
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

package com.fullteaching.backend.e2e;

import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeAll;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.TestInfo;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.runner.RunWith;
import org.junit.platform.runner.JUnitPlatform;
import org.junit.Assert;
import org.openqa.selenium.By;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.support.ui.ExpectedConditions;

import io.github.bonigarcia.SeleniumExtension;
import io.github.bonigarcia.wdm.ChromeDriverManager;
import io.github.bonigarcia.wdm.FirefoxDriverManager;

/**
 * E2E tests for FullTeaching video session.
 *
 * @author Pablo Fuente (pablo.fuente@urjc.es)
 */
@Tag("e2e")
@DisplayName("E2E tests for FullTeaching video session")
@ExtendWith(SeleniumExtension.class)
@RunWith(JUnitPlatform.class)
public class FullTeachingTestE2EDualVideoSession extends FullTeachingTestE2E {

	private static String TEACHER_BROWSER;
	private static String STUDENT_BROWSER;

	static Exception ex = null;

	final String teacherMail = "teacher@gmail.com";
	final String teacherPass = "pass";
	final String teacherName = "Teacher Cheater";
	final String studentMail = "student1@gmail.com";
	final String studentPass = "pass";
	final String studentName = "Student Imprudent";

	BrowserUser user;

	public FullTeachingTestE2EDualVideoSession() {
		super();
	}

	@BeforeAll()
	static void setupAll() {

		if (System.getenv("ET_EUS_API") == null) {
			// Outside ElasTest
			ChromeDriverManager.getInstance().setup();
			FirefoxDriverManager.getInstance().setup();
		}

		if (System.getenv("ET_SUT_HOST") != null) {
			APP_URL = "https://" + System.getenv("ET_SUT_HOST") + ":5000/";
		} else {
			APP_URL = System.getProperty("app.url");
			if (APP_URL == null) {
				APP_URL = "https://localhost:5000/";
			}
		}

		TEACHER_BROWSER = System.getenv("TEACHER_BROWSER");
		STUDENT_BROWSER = System.getenv("STUDENT_BROWSER");

		if ((TEACHER_BROWSER == null) || (!TEACHER_BROWSER.equals(FIREFOX))) {
			TEACHER_BROWSER = CHROME;
		}

		if ((STUDENT_BROWSER == null) || (!STUDENT_BROWSER.equals(FIREFOX))) {
			STUDENT_BROWSER = CHROME;
		}

		log.info("Using URL {} to connect to openvidu-testapp", APP_URL);
	}

	@AfterEach
	void dispose(TestInfo info) {
		try {
			this.logout(user);
			user.dispose();
		} finally {
			log.info("##### Finish test: " + info.getTestMethod().get().getName());
		}
	}

	@Test
	void oneToOneVideoAudioSessionChrome(TestInfo info) throws Exception {

		log.info("##### Start test: " + info.getTestMethod().get().getName());

		// TEACHER

		this.user = setupBrowser(TEACHER_BROWSER, info, "Teacher", 30);

		this.slowLogin(user, teacherMail, teacherPass);

		waitSeconds(1);

		log.info("{} entering first course", user.getClientData());

		user.getWaiter().until(ExpectedConditions.presenceOfElementLocated(
				By.cssSelector(("ul.collection li.collection-item:first-child div.course-title"))));
		user.getDriver().findElement(By.cssSelector("ul.collection li.collection-item:first-child div.course-title"))
				.click();

		waitSeconds(1);

		log.info("{} navigating to 'Sessions' tab", user.getClientData());

		user.getWaiter().until(ExpectedConditions.presenceOfElementLocated(By.cssSelector(("#md-tab-label-0-1"))));
		user.getDriver().findElement(By.cssSelector("#md-tab-label-0-1")).click();

		waitSeconds(1);

		log.info("{} getting into first session", user.getClientData());

		user.getDriver().findElement(By.cssSelector("ul div:first-child li.session-data div.session-ready")).click();

		waitSeconds(1);

		user.getWaiter().until(ExpectedConditions.presenceOfElementLocated(By.cssSelector(("div.participant video"))));

		waitSeconds(1);

		checkVideoPlaying(user, user.getDriver().findElement(By.cssSelector(("div.participant video"))),
				"div.participant");

		// STUDENT

		BrowserUser student = setupBrowser(STUDENT_BROWSER, info, "Student", 30);
		slowLogin(student, studentMail, studentPass);

		waitSeconds(1);

		log.info("{} entering first course", student.getClientData());

		student.getWaiter().until(ExpectedConditions.presenceOfElementLocated(
				By.cssSelector(("ul.collection li.collection-item:first-child div.course-title"))));
		student.getDriver().findElement(By.cssSelector("ul.collection li.collection-item:first-child div.course-title"))
				.click();

		waitSeconds(1);

		log.info("{} navigating to 'Courses' tab", student.getClientData());

		student.getWaiter().until(ExpectedConditions.presenceOfElementLocated(By.cssSelector(("#md-tab-label-0-1"))));
		student.getDriver().findElement(By.cssSelector("#md-tab-label-0-1")).click();

		waitSeconds(1);

		log.info("{} getting into first session", student.getClientData());

		student.getDriver().findElement(By.cssSelector("ul div:first-child li.session-data div.session-ready")).click();

		waitSeconds(1);

		student.getWaiter()
				.until(ExpectedConditions.presenceOfElementLocated(By.cssSelector(("div.participant video"))));

		waitSeconds(1);

		checkVideoPlaying(student, student.getDriver().findElement(By.cssSelector(("div.participant video"))),
				"div.participant");

		log.info("STARTING SIMPLEX SESSION");
		waitSeconds(30);
		log.info("ENDING SIMPLEX SESSION");

		// Student asks for intervention
		student.getWaiter().until(ExpectedConditions
				.elementToBeClickable(By.xpath("//div[@id='div-header-buttons']//i[text() = 'record_voice_over']")));

		log.info("{} asking for intervention", student.getClientData());

		student.getDriver().findElement(By.xpath("//div[@id='div-header-buttons']//i[text() = 'record_voice_over']"))
				.click();

		waitSeconds(1);

		// Teacher accepts intervention
		user.getWaiter().until(ExpectedConditions.elementToBeClickable(By.xpath("//a[contains(@class, 'usr-btn')]")));

		log.info("{} accepts student intervention", user.getClientData());

		user.getDriver().findElement(By.xpath("//a[contains(@class, 'usr-btn')]")).click();

		// Check both videos for both users
		student.getWaiter()
				.until(ExpectedConditions.presenceOfElementLocated(By.cssSelector(("div.participant-small video"))));
		// Small video of student
		checkVideoPlaying(student, student.getDriver().findElement(By.cssSelector(("div.participant-small video"))),
				"div.participant-small");
		// Main video of student
		checkVideoPlaying(student, student.getDriver().findElement(By.cssSelector(("div.participant video"))),
				"div.participant");

		user.getWaiter()
				.until(ExpectedConditions.presenceOfElementLocated(By.cssSelector(("div.participant-small video"))));
		// Small video of teacher
		checkVideoPlaying(user, user.getDriver().findElement(By.cssSelector(("div.participant-small video"))),
				"div.participant-small");
		// Main video of teacher
		checkVideoPlaying(user, user.getDriver().findElement(By.cssSelector(("div.participant video"))),
				"div.participant");

		log.info("STARTING DUPLEX SESSION");
		waitSeconds(30);
		log.info("ENDING DUPLEX SESSION");

		// Teacher stops student intervention
		user.getWaiter().until(ExpectedConditions.elementToBeClickable(By.xpath("//a[contains(@class, 'usr-btn')]")));

		log.info("{} canceling student intervention", user.getClientData());

		user.getDriver().findElement(By.xpath("//a[contains(@class, 'usr-btn')]")).click();

		// Wait until only one video
		user.getWaiter().until(ExpectedConditions.not(
				ExpectedConditions.presenceOfAllElementsLocatedBy(By.cssSelector(("div.participant-small video")))));
		student.getWaiter().until(ExpectedConditions.not(
				ExpectedConditions.presenceOfAllElementsLocatedBy(By.cssSelector(("div.participant-small video")))));

		waitSeconds(4);

		// Logout student
		this.logout(student);
		student.dispose();

	}

	private boolean checkVideoPlaying(BrowserUser user, WebElement videoElement, String containerQuerySelector) {

		log.info("{} waiting for video in container '{}' to be playing", user.getClientData(), containerQuerySelector);

		// Video element should be in 'readyState'='HAVE_ENOUGH_DATA'
		user.getWaiter().until(ExpectedConditions.attributeToBe(videoElement, "readyState", "4"));

		// Video should have a srcObject (type MediaStream) with the attribute 'active'
		// to true
		Assert.assertTrue((boolean) user.runJavascript("return document.querySelector('" + containerQuerySelector
				+ "').getElementsByTagName('video')[0].srcObject.active"));

		// Video should trigger 'playing' event
		user.runJavascript("document.querySelector('" + containerQuerySelector
				+ "').getElementsByTagName('video')[0].addEventListener('playing', window.MY_FUNC('"
				+ containerQuerySelector + "'));");

		user.getWaiter()
				.until(ExpectedConditions.attributeContains(By.id("video-playing-div"), "innerHTML", "VIDEO PLAYING"));

		user.runJavascript("document.body.removeChild(document.getElementById('video-playing-div'))");

		return true;
	}

}
