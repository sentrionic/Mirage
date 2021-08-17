package xyz.mirage.app.presentation.navigation

sealed class Graph(
    val route: String,
) {
    // Auth
    object Auth : Graph("auth-graph")
    object Home : Graph("home-graph")
    object Search : Graph("search-graph")
    object Account : Graph("account-graph")
}