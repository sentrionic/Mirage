package xyz.mirage.app.business.interactors.post

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.flow
import xyz.mirage.app.business.datasources.cache.post.PostDao
import xyz.mirage.app.business.datasources.cache.post.toEntity
import xyz.mirage.app.business.datasources.cache.post.toPostList
import xyz.mirage.app.business.datasources.network.core.handleUseCaseException
import xyz.mirage.app.business.datasources.network.main.post.PostService
import xyz.mirage.app.business.datasources.network.main.post.dto.toPostList
import xyz.mirage.app.business.domain.core.DataState
import xyz.mirage.app.business.domain.models.Post

class GetFeed(
    private val cache: PostDao,
    private val service: PostService,
) {

    fun execute(
        page: Int,
        cursor: String?,
        isNetworkAvailable: Boolean,
    ): Flow<DataState<List<Post>>> = flow {
        emit(DataState.loading())

        if (isNetworkAvailable) {
            val posts = service.feed(cursor = cursor).posts.toPostList()
            // insert into cache

            for (post in posts) {
                cache.insertPost(post.toEntity())
                cache.insertAuthor(post.profile.toEntity())
            }
        }

        // query the cache
        val cacheResult = cache.getFeed(
            page = page
        )

        // emit List<Recipe> from cache
        val list = cacheResult.toPostList()

        emit(DataState.data(response = null, data = list))
    }.catch { e ->
        emit(handleUseCaseException(e))
    }
}