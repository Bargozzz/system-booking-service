import { defineComponent, ref, mergeProps, unref, useSSRContext } from 'vue';
import { ssrRenderAttrs, ssrIncludeBooleanAttr, ssrLooseContain, ssrLooseEqual, ssrRenderAttr, ssrInterpolate, ssrRenderList, ssrRenderClass } from 'vue/server-renderer';
import { u as useBookingStore } from './booking-DjwprT6X.mjs';
import { _ as _export_sfc, u as useApi, b as useRouter } from './server.mjs';
import 'pinia';
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
import 'vue-router';

const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "index",
  __ssrInlineRender: true,
  setup(__props) {
    const bookingStore = useBookingStore();
    useApi();
    useRouter();
    const loading = ref(false);
    const searched = ref(false);
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<div${ssrRenderAttrs(mergeProps({ class: "search-page" }, _attrs))} data-v-f51d6e16><section class="hero-section" data-v-f51d6e16><div class="hero-content" data-v-f51d6e16><h1 data-v-f51d6e16>Pesan Tiket Bus Real-Time</h1><p data-v-f51d6e16>Sistem reservasi instan dengan fitur penguncian kursi transaksional 5 menit.</p></div><div class="search-card glass-card" data-v-f51d6e16><form class="search-form" data-v-f51d6e16><div class="form-group" data-v-f51d6e16><label data-v-f51d6e16>Kota Asal</label><select class="input-field" data-v-f51d6e16><option value="Jakarta" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchOrigin) ? ssrLooseContain(unref(bookingStore).searchOrigin, "Jakarta") : ssrLooseEqual(unref(bookingStore).searchOrigin, "Jakarta")) ? " selected" : ""}>Jakarta</option><option value="Bandung" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchOrigin) ? ssrLooseContain(unref(bookingStore).searchOrigin, "Bandung") : ssrLooseEqual(unref(bookingStore).searchOrigin, "Bandung")) ? " selected" : ""}>Bandung</option><option value="Surabaya" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchOrigin) ? ssrLooseContain(unref(bookingStore).searchOrigin, "Surabaya") : ssrLooseEqual(unref(bookingStore).searchOrigin, "Surabaya")) ? " selected" : ""}>Surabaya</option><option value="Malang" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchOrigin) ? ssrLooseContain(unref(bookingStore).searchOrigin, "Malang") : ssrLooseEqual(unref(bookingStore).searchOrigin, "Malang")) ? " selected" : ""}>Malang</option></select></div><div class="form-group" data-v-f51d6e16><label data-v-f51d6e16>Kota Tujuan</label><select class="input-field" data-v-f51d6e16><option value="Bandung" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchDestination) ? ssrLooseContain(unref(bookingStore).searchDestination, "Bandung") : ssrLooseEqual(unref(bookingStore).searchDestination, "Bandung")) ? " selected" : ""}>Bandung</option><option value="Jakarta" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchDestination) ? ssrLooseContain(unref(bookingStore).searchDestination, "Jakarta") : ssrLooseEqual(unref(bookingStore).searchDestination, "Jakarta")) ? " selected" : ""}>Jakarta</option><option value="Malang" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchDestination) ? ssrLooseContain(unref(bookingStore).searchDestination, "Malang") : ssrLooseEqual(unref(bookingStore).searchDestination, "Malang")) ? " selected" : ""}>Malang</option><option value="Surabaya" data-v-f51d6e16${ssrIncludeBooleanAttr(Array.isArray(unref(bookingStore).searchDestination) ? ssrLooseContain(unref(bookingStore).searchDestination, "Surabaya") : ssrLooseEqual(unref(bookingStore).searchDestination, "Surabaya")) ? " selected" : ""}>Surabaya</option></select></div><div class="form-group" data-v-f51d6e16><label data-v-f51d6e16>Tanggal Keberangkatan</label><input${ssrRenderAttr("value", unref(bookingStore).searchDate)} type="date" class="input-field" required data-v-f51d6e16></div><button type="submit" class="btn-primary btn-search"${ssrIncludeBooleanAttr(unref(loading)) ? " disabled" : ""} data-v-f51d6e16>`);
      if (unref(loading)) {
        _push(`<span data-v-f51d6e16>Mencari...</span>`);
      } else {
        _push(`<span data-v-f51d6e16>\u{1F50D} Cari Jadwal</span>`);
      }
      _push(`</button></form></div></section><section class="results-section" data-v-f51d6e16>`);
      if (unref(searched)) {
        _push(`<div class="section-header" data-v-f51d6e16><h2 data-v-f51d6e16>Hasil Pencarian (${ssrInterpolate(unref(bookingStore).schedules.length)})</h2><p data-v-f51d6e16>Menampilkan jadwal dari ${ssrInterpolate(unref(bookingStore).searchOrigin)} ke ${ssrInterpolate(unref(bookingStore).searchDestination)}</p></div>`);
      } else {
        _push(`<!---->`);
      }
      if (unref(loading)) {
        _push(`<div class="loading-state glass-card" data-v-f51d6e16><p data-v-f51d6e16>Memuat data jadwal bus...</p></div>`);
      } else if (unref(searched) && unref(bookingStore).schedules.length === 0) {
        _push(`<div class="empty-state glass-card" data-v-f51d6e16><span class="empty-icon" data-v-f51d6e16>\u{1F6AB}</span><h3 data-v-f51d6e16>Tidak Ada Jadwal Ditemukan</h3><p data-v-f51d6e16>Coba pilih kombinasi kota asal/tujuan atau tanggal keberangkatan yang lain.</p></div>`);
      } else {
        _push(`<div class="schedule-grid" data-v-f51d6e16><!--[-->`);
        ssrRenderList(unref(bookingStore).schedules, (schedule) => {
          _push(`<div class="schedule-card glass-card" data-v-f51d6e16><div class="operator-info" data-v-f51d6e16><div class="operator-badge" data-v-f51d6e16><span class="operator-code" data-v-f51d6e16>${ssrInterpolate(schedule.operator_code)}</span><span class="operator-name" data-v-f51d6e16>${ssrInterpolate(schedule.operator)}</span></div><div class="${ssrRenderClass([schedule.available_seats ? "badge-success" : "badge-warning", "seats-badge badge"])}" data-v-f51d6e16>${ssrInterpolate(schedule.available_seats)} dari ${ssrInterpolate(schedule.total_seats)} kursi tersedia </div></div><div class="route-details" data-v-f51d6e16><div class="time-block" data-v-f51d6e16><span class="time" data-v-f51d6e16>${ssrInterpolate(schedule.departure_time.substring(0, 5))}</span><span class="city" data-v-f51d6e16>${ssrInterpolate(schedule.origin)}</span></div><div class="route-line" data-v-f51d6e16><span class="duration" data-v-f51d6e16>Langsung</span><div class="line" data-v-f51d6e16></div></div><div class="time-block" data-v-f51d6e16><span class="time" data-v-f51d6e16>Tiba</span><span class="city" data-v-f51d6e16>${ssrInterpolate(schedule.destination)}</span></div></div><div class="card-footer" data-v-f51d6e16><div class="price-tag" data-v-f51d6e16><span class="price-label" data-v-f51d6e16>Harga per kursi</span><span class="price-val" data-v-f51d6e16>Rp ${ssrInterpolate(schedule.price.toLocaleString("id-ID"))}</span></div><button class="btn-primary"${ssrIncludeBooleanAttr(!schedule.available_seats) ? " disabled" : ""} data-v-f51d6e16> Pilih Kursi \u2794 </button></div></div>`);
        });
        _push(`<!--]--></div>`);
      }
      _push(`</section></div>`);
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/index.vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
const index = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-f51d6e16"]]);

export { index as default };
//# sourceMappingURL=index-CuCtzVDD.mjs.map
