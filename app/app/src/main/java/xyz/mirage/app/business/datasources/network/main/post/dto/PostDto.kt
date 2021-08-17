package xyz.mirage.app.business.datasources.network.main.post.dto

import com.squareup.moshi.Json
import xyz.mirage.app.business.domain.models.Post

data class PostDto(
    @Json(name = "id")
    val id: String,

    @Json(name = "text")
    val text: String?,

    @Json(name = "likes")
    val likes: Int = 0,

    @Json(name = "liked")
    val liked: Boolean,

    @Json(name = "retweets")
    val retweets: Int = 0,

    @Json(name = "retweeted")
    val retweeted: Boolean,

    @Json(name = "isRetweet")
    val isRetweet: Boolean,

    @Json(name = "file")
    val file: AttachmentDto?,

    @Json(name = "author")
    val author: ProfileDto,

    @Json(name = "createdAt")
    val createdAt: String
) {
    fun toPost(): Post {
        return Post(
            id = id,
            text = text,
            likes = likes,
            liked = liked,
            retweets = retweets,
            retweeted = retweeted,
            isRetweet = isRetweet,
            file = file?.toAttachment(),
            profile = author.toProfile(),
            createdAt = createdAt
        )
    }
}

fun List<PostDto>.toPostList(): List<Post> {
    return map { it.toPost() }
}
