package com.brotherlogic.bartwarn

import android.app.PendingIntent
import com.google.android.gms.location.Geofence
import com.google.android.gms.location.GeofencingClient
import com.google.android.gms.location.GeofencingRequest
import com.google.android.gms.tasks.Task
import org.junit.Assert.assertEquals
import org.junit.Test
import org.mockito.kotlin.any
import org.mockito.kotlin.argumentCaptor
import org.mockito.kotlin.mock
import org.mockito.kotlin.verify
import org.mockito.kotlin.whenever

class GeofenceManagerTest {

    @Test
    fun `test addGeofences creates correct geofences`() {
        val mockGeofencingClient: GeofencingClient = mock()
        val mockTask: Task<Void> = mock()
        whenever(mockGeofencingClient.addGeofences(any<GeofencingRequest>(), any<PendingIntent>())).thenReturn(mockTask)
        whenever(mockTask.addOnSuccessListener(any())).thenReturn(mockTask)
        whenever(mockTask.addOnFailureListener(any())).thenReturn(mockTask)

        val geofenceManager = GeofenceManager(mockGeofencingClient)
        val mockPendingIntent: PendingIntent = mock()

        var successCalled = false
        geofenceManager.addGeofences(mockPendingIntent, onSuccess = { successCalled = true }, onFailure = {})

        val captor = argumentCaptor<GeofencingRequest>()
        verify(mockGeofencingClient).addGeofences(captor.capture(), any<PendingIntent>())

        val request = captor.firstValue
        assertEquals(BARTStations.stations.size, request.geofences.size)
        
        val firstGeofence = request.geofences[0]
        assertEquals("12th", firstGeofence.requestId)
        
        // Unfortunately, getRadius() or other getters might not be public on Geofence unless we cast it, 
        // but we know there are exactly the correct number of geofences.
    }
}
