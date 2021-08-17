package xyz.mirage.app.presentation.ui.main.home.list.components

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import xyz.mirage.app.business.domain.models.Post

@Composable
fun PostActions(
    post: Post,
    onToggleLike: (String) -> Unit,
    onToggleRetweet: (String) -> Unit,
) {
    Row(
        modifier = Modifier
            .fillMaxWidth()
            .padding(end = 40.dp),
        verticalAlignment = Alignment.CenterVertically,
        horizontalArrangement = Arrangement.spacedBy(30.dp)
    ) {

        Row(
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            AnimatedRetweetButton(
                modifier = Modifier,
                isRetweeted = post.retweeted,
                onToggle = {
                    onToggleRetweet(post.id)
                }
            )

            Text(text = post.retweets.toString())
        }

        Row(
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.spacedBy(4.dp)
        ) {
            AnimatedHeartButton(
                modifier = Modifier,
                isLiked = post.liked,
                onToggle = {
                    onToggleLike(post.id)
                },
            )

            Text(text = post.likes.toString())
        }
    }
}
