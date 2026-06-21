package com.brotherlogic.bartwarn

import android.content.Context
import androidx.test.core.app.ApplicationProvider
import androidx.work.ListenableWorker
import androidx.work.testing.TestListenableWorkerBuilder
import kotlinx.coroutines.runBlocking
import org.junit.Assert.assertEquals
import org.junit.Test
import org.junit.runner.RunWith
import org.robolectric.RobolectricTestRunner

@RunWith(RobolectricTestRunner::class)
class LocationPingWorkerTest {
    @Test
    fun testWorkerFailsWithoutStationId() {
        val context = ApplicationProvider.getApplicationContext<Context>()
        val worker = TestListenableWorkerBuilder<LocationPingWorker>(context).build()
        runBlocking {
            val result = worker.doWork()
            assertEquals(ListenableWorker.Result.failure(), result)
        }
    }
}
