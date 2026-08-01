import { createRouter, createWebHistory } from "vue-router";
import { authState } from "../store/auth";
import Login from "../views/Login.vue";
import Search from "../views/Search.vue";
import SeatSelection from "../views/SeatSelection.vue";
import BookingSummary from "../views/BookingSummary.vue";
import Confirmation from "../views/Confirmation.vue";

const routes = [
  { path: "/", redirect: "/search" },
  { path: "/login", name: "Login", component: Login },
  { path: "/search", name: "Search", component: Search, meta: { requiresAuth: false } },
  {
    path: "/schedules/:scheduleId/seats",
    name: "SeatSelection",
    component: SeatSelection,
    meta: { requiresAuth: true },
    props: true,
  },
  {
    path: "/booking-summary",
    name: "BookingSummary",
    component: BookingSummary,
    meta: { requiresAuth: true },
  },
  {
    path: "/confirmation",
    name: "Confirmation",
    component: Confirmation,
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// Seat locking/booking requires login; searching schedules does not.
router.beforeEach((to) => {
  if (to.meta.requiresAuth && !authState.isAuthenticated) {
    return { name: "Login", query: { redirect: to.fullPath } };
  }
});

export default router;
