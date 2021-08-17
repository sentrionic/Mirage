package xyz.mirage.app.business.datasources.cache.post

import androidx.room.Embedded
import androidx.room.Relation
import com.squareup.moshi.JsonAdapter
import com.squareup.moshi.Moshi
import xyz.mirage.app.business.domain.models.Attachment
import xyz.mirage.app.business.domain.models.Post
import xyz.mirage.app.presentation.core.util.DateUtils

val moshi: Moshi = Moshi.Builder().build()
val fileAdapter: JsonAdapter<Attachment> = moshi.adapter(Attachment::class.java)

data class PostAuthor(
    @Embedded
    val post: PostEntity,
    @Relation(
        parentColumn = "authorId",
        entityColumn = "id"
    )
    val author: ProfileEntity
) {
    fun toPost(): Post {
        return Post(
            id = post.id,
            text = post.text,
            likes = post.likes,
            liked = post.liked,
            retweets = post.retweets,
            retweeted = post.retweeted,
            isRetweet = post.isRetweet,
            file = post.file?.let { convertStringToFile(it) },
            profile = author.toProfile(),
            createdAt = post.createdAt,
        )
    }

    private fun convertStringToFile(file: String): Attachment {
        return fileAdapter.fromJson(file) ?: throw Exception("Couldn't restore attachment")
    }

    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as PostAuthor

        if (post != other.post) return false
        if (author != other.author) return false

        return true
    }

    override fun hashCode(): Int {
        var result = post.hashCode()
        result = 31 * result + author.hashCode()
        return result
    }
}

private fun convertFileToString(file: Attachment): String {
    return fileAdapter.toJson(file)
}

fun Post.toEntity(isFeed: Boolean = true): PostEntity {
    return PostEntity(
        id = id,
        text = text,
        likes = likes,
        liked = liked,
        retweets = retweets,
        retweeted = retweeted,
        isRetweet = isRetweet,
        file = file?.let { convertFileToString(it) },
        authorId = profile.id,
        createdAt = createdAt,
        dateCached = DateUtils.createTimestamp(),
        isFeed = isFeed
    )
}

fun List<PostAuthor>.toPostList(): List<Post> {
    return map { it.toPost() }
}