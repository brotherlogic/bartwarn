package com.brotherlogic.bartwarn

import android.content.Context
import android.content.Intent
import org.junit.Test
import org.mockito.kotlin.mock
import org.mockito.kotlin.verifyNoInteractions

class BootCompletedReceiverTest {

    @Test
    fun `onReceive does nothing if action is not boot completed`() {
        val receiver = BootCompletedReceiver()
        val mockContext: Context = mock()
        val mockIntent: Intent = mock()
        
        // This will throw exception if we didn't mock action
        // Actually mockIntent.action will return null by default
        
        receiver.onReceive(mockContext, mockIntent)
        
        verifyNoInteractions(mockContext)
    }
}
