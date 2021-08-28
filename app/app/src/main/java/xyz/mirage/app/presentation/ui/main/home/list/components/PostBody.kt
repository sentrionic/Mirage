package xyz.mirage.app.presentation.ui.main.home.list.components

import androidx.compose.foundation.Image
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.Box
import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.draw.clip
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.layout.ContentScale
import androidx.compose.ui.unit.dp
import coil.ImageLoader
import coil.annotation.ExperimentalCoilApi
import coil.compose.ImagePainter
import coil.compose.rememberImagePainter
import xyz.mirage.app.business.domain.models.Post

@ExperimentalCoilApi
@Composable
fun PostBody(
    post: Post,
    imageLoader: ImageLoader,
) {
    post.text?.let {
        Text(
            text = it,
        )
    }

    post.file?.let { file ->

        val painter = rememberImagePainter(
            data = file.url,
            imageLoader = imageLoader,
        )

        Spacer(modifier = Modifier.height(10.dp))

        Image(
            painter = painter,
            contentDescription = "Post Image for ${post.id}",
            modifier = Modifier
                .height(180.dp)
                .fillMaxWidth()
                .clip(shape = RoundedCornerShape(10.dp)),
            contentScale = ContentScale.Crop
        )

        when (painter.state) {
            is ImagePainter.State.Loading -> {
                Box(
                    modifier = Modifier
                        .background(Color.Transparent)
                        .height(180.dp)
                        .fillMaxWidth()
                        .clip(shape = RoundedCornerShape(10.dp)),
                )
            }
            else -> {
            }
        }
    }
}