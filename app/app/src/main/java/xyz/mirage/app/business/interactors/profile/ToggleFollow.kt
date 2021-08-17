package xyz.mirage.app.business.interactors.profile

import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.catch
import kotlinx.coroutines.flow.flow
import xyz.mirage.app.business.datasources.cache.profile.ProfileDao
import xyz.mirage.app.business.datasources.network.core.handleUseCaseException
import xyz.mirage.app.business.datasources.network.main.profile.ProfileService
import xyz.mirage.app.business.domain.core.*
import xyz.mirage.app.business.domain.models.Profile

class ToggleFollow(
    private val service: ProfileService,
    private val cache: ProfileDao,
) {

    fun execute(
        username: String
    ): Flow<DataState<Profile>> = flow {
        emit(DataState.loading())

        val response = service.toggleFollow(username)

        // Update cache
        cache.updateFollow(
            followers = response.followers,
            followee = response.followee,
            following = response.following,
            id = response.id
        )

        emit(
            DataState.data(
                data = response.toProfile(),
                response = Response(
                    message = SuccessHandling.SUCCESS_FOLLOW_TOGGLED,
                    uiComponentType = UIComponentType.None(),
                    messageType = MessageType.Success()
                )
            )
        )
    }.catch { e ->
        emit(handleUseCaseException(e))
    }
}