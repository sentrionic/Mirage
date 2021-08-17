package xyz.mirage.app.presentation.ui.main.profile.components

import androidx.compose.runtime.Composable
import xyz.mirage.app.presentation.ui.shared.CustomDivider

private const val PAGE_SIZE = 20

@Composable
fun ListItem(
    index: Int,
    page: Int,
    isLoading: Boolean,
    isDarkTheme: Boolean,
    onChangeScrollPosition: () -> Unit,
    fetchNextPage: () -> Unit,
    content: @Composable () -> Unit
) {
    onChangeScrollPosition()

    if ((index + 1) >= (page * PAGE_SIZE) && !isLoading) {
        fetchNextPage()
    }

    content()
    CustomDivider(isDarkTheme = isDarkTheme)
}