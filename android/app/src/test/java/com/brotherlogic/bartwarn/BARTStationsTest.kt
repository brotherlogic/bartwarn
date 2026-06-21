package com.brotherlogic.bartwarn

import org.junit.Test
import org.junit.Assert.*

class BARTStationsTest {
    @Test
    fun testBARTStationsExistsAndHasData() {
        val stations = BARTStations.stations
        assertNotNull(stations)
        assertTrue(stations.isNotEmpty())
        
        val firstStation = stations.first()
        assertNotNull(firstStation.id)
        assertNotNull(firstStation.name)
        assertNotNull(firstStation.lat)
        assertNotNull(firstStation.lng)
    }
}
