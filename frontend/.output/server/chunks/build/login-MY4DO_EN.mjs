import { defineComponent, ref, mergeProps, unref, useSSRContext } from 'vue';
import { ssrRenderAttrs, ssrRenderAttr, ssrInterpolate, ssrIncludeBooleanAttr, ssrRenderList } from 'vue/server-renderer';
import { _ as _export_sfc, c as useAuthStore, u as useApi, b as useRouter } from './server.mjs';
import '../_/nitro.mjs';
import 'node:http';
import 'node:https';
import 'node:events';
import 'node:buffer';
import 'node:fs';
import 'node:path';
import 'node:crypto';
import 'node:url';
import '../routes/renderer.mjs';
import 'vue-bundle-renderer/runtime';
import 'unhead/server';
import 'devalue';
import 'unhead/utils';
import 'unhead/plugins';
import 'pinia';
import 'vue-router';

const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "login",
  __ssrInlineRender: true,
  setup(__props) {
    const email = ref("user1@example.com");
    const password = ref("password");
    const loading = ref(false);
    const error = ref(null);
    useAuthStore();
    useApi();
    useRouter();
    const demoUsers = [
      { name: "User Satu", email: "user1@example.com" },
      { name: "User Dua", email: "user2@example.com" },
      { name: "User Tiga", email: "user3@example.com" }
    ];
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps({ class: "login-page" }, _attrs))} data-v-315f10cc><div class="login-card glass-card" data-v-315f10cc><div class="header" data-v-315f10cc><span class="logo" data-v-315f10cc>\u{1F68C}</span><h2 data-v-315f10cc>Login MiniBooking</h2><p data-v-315f10cc>Masuk ke akun untuk mulai mengunci dan memesan kursi bus</p></div><form class="login-form" data-v-315f10cc><div class="form-group" data-v-315f10cc><label data-v-315f10cc>Email Address</label><input${ssrRenderAttr("value", unref(email))} type="email" class="input-field" placeholder="user1@example.com" required data-v-315f10cc></div><div class="form-group" data-v-315f10cc><label data-v-315f10cc>Password</label><input${ssrRenderAttr("value", unref(password))} type="password" class="input-field" placeholder="\u2022\u2022\u2022\u2022\u2022\u2022\u2022\u2022" required data-v-315f10cc></div>`);
      if (unref(error)) {
        _push(`<div class="error-banner" data-v-315f10cc>${ssrInterpolate(unref(error))}</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`<button type="submit" class="btn-primary btn-block"${ssrIncludeBooleanAttr(unref(loading)) ? " disabled" : ""} data-v-315f10cc>`);
      if (unref(loading)) {
        _push(`<span data-v-315f10cc>Loading...</span>`);
      } else {
        _push(`<span data-v-315f10cc>Login Sekarang</span>`);
      }
      _push(`</button><div class="quick-logins" data-v-315f10cc><span data-v-315f10cc>Akun Demo Uji Coba:</span><div class="demo-buttons" data-v-315f10cc><!--[-->`);
      ssrRenderList(demoUsers, (u) => {
        _push(`<button type="button" class="btn-demo-pill" data-v-315f10cc>${ssrInterpolate(u.name)} (${ssrInterpolate(u.email)}) </button>`);
      });
      _push(`<!--]--></div></div></form></div></div>`);
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/login.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
const login = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-315f10cc"]]);

export { login as default };
//# sourceMappingURL=login-MY4DO_EN.mjs.map
