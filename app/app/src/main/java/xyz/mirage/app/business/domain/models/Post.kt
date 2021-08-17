package xyz.mirage.app.business.domain.models

data class Post(
    val id: String,
    val text: String?,
    var likes: Int = 0,
    var liked: Boolean,
    var retweets: Int = 0,
    var retweeted: Boolean,
    val isRetweet: Boolean,
    val file: Attachment?,
    val profile: Profile,
    val createdAt: String
)