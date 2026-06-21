package com.brotherlogic.bartwarn

import android.app.PendingIntent
import org.junit.Test
import org.mockito.kotlin.any
import org.mockito.kotlin.mock
import org.mockito.kotlin.verify

class BootCompletedReceiverTest {

    @Test
    fun `test handleBootCompleted registers geofences`() {
        val receiver = BootCompletedReceiver()
        val mockGeofenceManager: GeofenceManager = mock()
        val mockPendingIntent: PendingIntent = mock()

        receiver.handleBootCompleted(mockGeofenceManager, mockPendingIntent)
        
        verify(mockGeofenceManager).addGeofences(any(), any(), any())
    }
}
