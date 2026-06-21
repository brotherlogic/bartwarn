package com.brotherlogic.bartwarn

import android.content.Context
import android.content.Intent
import org.junit.Test
import org.mockito.kotlin.mock
import org.mockito.kotlin.whenever

class BootCompletedReceiverTest {

    @Test
    fun `test onReceive handles ACTION_BOOT_COMPLETED`() {
        val receiver = BootCompletedReceiver()
        val context: Context = mock()
        val intent: Intent = mock()
        whenever(intent.action).thenReturn(Intent.ACTION_BOOT_COMPLETED)

        receiver.onReceive(context, intent)
    }
}
