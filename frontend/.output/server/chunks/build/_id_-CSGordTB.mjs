import { _ as _export_sfc, d as useRoute, b as useRouter, u as useApi, c as useAuthStore, a as __nuxt_component_0$1 } from './server.mjs';
import { defineComponent, computed, ref, mergeProps, withCtx, createTextVNode, unref, useSSRContext } from 'vue';
import { ssrRenderAttrs, ssrRenderComponent, ssrInterpolate, ssrRenderClass, ssrIncludeBooleanAttr, ssrRenderList } from 'vue/server-renderer';
import { u as useBookingStore } from './booking-DjwprT6X.mjs';
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

const _sfc_main$3 = /* @__PURE__ */ defineComponent({
  __name: "CountdownTimer",
  __ssrInlineRender: true,
  props: {
    seconds: {}
  },
  setup(__props) {
    const props = __props;
    const expired = computed(() => props.seconds <= 0);
    const formattedTime = computed(() => {
      const m = Math.floor(props.seconds / 60);
      const s = props.seconds % 60;
      return `${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
    });
    return (_ctx, _push, _parent, _attrs) => {
      if (__props.seconds > 0) {
        _push(`<div${ssrRenderAttrs(mergeProps({ class: "timer-card pulse-lock" }, _attrs))} data-v-24e41eb4><span class="timer-icon" data-v-24e41eb4>\u23F1\uFE0F</span><div class="timer-content" data-v-24e41eb4><span class="timer-title" data-v-24e41eb4>Kursi Dikunci Sementara</span><span class="timer-clock" data-v-24e41eb4>${ssrInterpolate(unref(formattedTime))}</span></div></div>`);
      } else if (unref(expired)) {
        _push(`<div${ssrRenderAttrs(mergeProps({ class: "timer-card timer-expired" }, _attrs))} data-v-24e41eb4><span class="timer-icon" data-v-24e41eb4>\u26A0\uFE0F</span><span class="timer-title" data-v-24e41eb4>Waktu kunci kursi habis! Silakan pilih kembali.</span></div>`);
      } else {
        _push(`<!---->`);
      }
    };
  }
});
const _sfc_setup$3 = _sfc_main$3.setup;
_sfc_main$3.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("components/CountdownTimer.vue");
  return _sfc_setup$3 ? _sfc_setup$3(props, ctx) : void 0;
};
const __nuxt_component_1 = /* @__PURE__ */ _export_sfc(_sfc_main$3, [["__scopeId", "data-v-24e41eb4"]]);
const _sfc_main$2 = /* @__PURE__ */ defineComponent({
  __name: "SeatItem",
  __ssrInlineRender: true,
  props: {
    seat: {},
    isSelected: { type: Boolean },
    currentUserId: {}
  },
  emits: ["select"],
  setup(__props) {
    const props = __props;
    const isLockedByMe = computed(() => {
      return props.isSelected;
    });
    const statusClass = computed(() => {
      if (props.isSelected) return "status-locked-user pulse-lock";
      if (props.seat.status === "booked") return "status-booked";
      if (props.seat.status === "locked") return "status-locked-other";
      return "status-available";
    });
    const statusText = computed(() => {
      if (props.isSelected) return "Dikunci Anda";
      if (props.seat.status === "booked") return "Terisi";
      if (props.seat.status === "locked") return "Dikunci User Lain";
      return "Tersedia";
    });
    return (_ctx, _push, _parent, _attrs) => {
      _push(`<button${ssrRenderAttrs(mergeProps({
        class: ["seat-btn", [
          unref(statusClass),
          { "is-selected": __props.isSelected }
        ]],
        disabled: _ctx.status === "booked" || _ctx.status === "locked" && !unref(isLockedByMe)
      }, _attrs))} data-v-7415fbb7><span class="seat-num" data-v-7415fbb7>${ssrInterpolate(__props.seat.seat_number)}</span><span class="seat-status-label" data-v-7415fbb7>${ssrInterpolate(unref(statusText))}</span></button>`);
    };
  }
});
const _sfc_setup$2 = _sfc_main$2.setup;
_sfc_main$2.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("components/SeatItem.vue");
  return _sfc_setup$2 ? _sfc_setup$2(props, ctx) : void 0;
};
const __nuxt_component_0 = /* @__PURE__ */ _export_sfc(_sfc_main$2, [["__scopeId", "data-v-7415fbb7"]]);
const _sfc_main$1 = /* @__PURE__ */ defineComponent({
  __name: "SeatGrid",
  __ssrInlineRender: true,
  props: {
    seats: {},
    selectedSeatId: {}
  },
  emits: ["select-seat"],
  setup(__props, { emit: __emit }) {
    const props = __props;
    const emit = __emit;
    const rows = ["A", "B", "C", "D", "E"];
    const fallbackSeat = (code) => ({
      id: 0,
      seat_number: code,
      status: "available",
      locked_until: null
    });
    const getSeat = (code) => {
      return props.seats.find((s) => s.seat_number === code) || fallbackSeat(code);
    };
    const onSelectSeat = (seat) => {
      if (seat.id !== 0) {
        emit("select-seat", seat);
      }
    };
    return (_ctx, _push, _parent, _attrs) => {
      const _component_SeatItem = __nuxt_component_0;
      _push(`<div${ssrRenderAttrs(mergeProps({ class: "bus-container glass-card" }, _attrs))} data-v-42f0aa19><div class="bus-front" data-v-42f0aa19><span class="driver-wheel" data-v-42f0aa19>\u2638\uFE0F Supir Bus</span></div><div class="seat-grid" data-v-42f0aa19><!--[-->`);
      ssrRenderList(rows, (row) => {
        _push(`<div class="grid-row" data-v-42f0aa19><div class="seat-pair" data-v-42f0aa19><!--[-->`);
        ssrRenderList([1, 2], (col) => {
          _push(ssrRenderComponent(_component_SeatItem, {
            key: `${row}${col}`,
            seat: getSeat(`${row}${col}`),
            "is-selected": __props.selectedSeatId === getSeat(`${row}${col}`).id,
            onSelect: onSelectSeat
          }, null, _parent));
        });
        _push(`<!--]--></div><div class="aisle" data-v-42f0aa19><span class="row-label" data-v-42f0aa19>${ssrInterpolate(row)}</span></div><div class="seat-pair" data-v-42f0aa19><!--[-->`);
        ssrRenderList([3, 4], (col) => {
          _push(ssrRenderComponent(_component_SeatItem, {
            key: `${row}${col}`,
            seat: getSeat(`${row}${col}`),
            "is-selected": __props.selectedSeatId === getSeat(`${row}${col}`).id,
            onSelect: onSelectSeat
          }, null, _parent));
        });
        _push(`<!--]--></div></div>`);
      });
      _push(`<!--]--></div><div class="seat-legend" data-v-42f0aa19><div class="legend-item" data-v-42f0aa19><span class="legend-dot status-available" data-v-42f0aa19></span> Tersedia</div><div class="legend-item" data-v-42f0aa19><span class="legend-dot status-locked-user" data-v-42f0aa19></span> Dikunci Anda</div><div class="legend-item" data-v-42f0aa19><span class="legend-dot status-locked-other" data-v-42f0aa19></span> Dikunci User Lain</div><div class="legend-item" data-v-42f0aa19><span class="legend-dot status-booked" data-v-42f0aa19></span> Terisi</div></div></div>`);
    };
  }
});
const _sfc_setup$1 = _sfc_main$1.setup;
_sfc_main$1.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("components/SeatGrid.vue");
  return _sfc_setup$1 ? _sfc_setup$1(props, ctx) : void 0;
};
const __nuxt_component_2 = /* @__PURE__ */ _export_sfc(_sfc_main$1, [["__scopeId", "data-v-42f0aa19"]]);
const _sfc_main = /* @__PURE__ */ defineComponent({
  __name: "[id]",
  __ssrInlineRender: true,
  setup(__props) {
    const route = useRoute();
    useRouter();
    useApi();
    useAuthStore();
    const bookingStore = useBookingStore();
    route.params.id;
    const schedule = computed(() => bookingStore.selectedSchedule);
    const selectedSeat = ref(null);
    const actionLoading = ref(false);
    const message = ref(null);
    const messageType = ref("success");
    const confirmedBooking = ref(null);
    const isLockedByMe = computed(() => {
      if (!selectedSeat.value || !bookingStore.activeLockedSeat) return false;
      return selectedSeat.value.id === bookingStore.activeLockedSeat.id && bookingStore.lockRemainingSeconds > 0;
    });
    const handleSelectSeat = (seat) => {
      var _a;
      if (seat.status === "booked" || seat.status === "locked" && seat.id !== ((_a = bookingStore.activeLockedSeat) == null ? void 0 : _a.id)) {
        return;
      }
      selectedSeat.value = seat;
      message.value = null;
    };
    return (_ctx, _push, _parent, _attrs) => {
      var _a, _b;
      const _component_NuxtLink = __nuxt_component_0$1;
      const _component_CountdownTimer = __nuxt_component_1;
      const _component_SeatGrid = __nuxt_component_2;
      _push(`<div${ssrRenderAttrs(mergeProps({ class: "seat-selection-page" }, _attrs))} data-v-8c23d228><div class="top-bar" data-v-8c23d228>`);
      _push(ssrRenderComponent(_component_NuxtLink, {
        to: "/",
        class: "back-link"
      }, {
        default: withCtx((_, _push2, _parent2, _scopeId) => {
          if (_push2) {
            _push2(` \u2190 Kembali ke Cari Jadwal `);
          } else {
            return [
              createTextVNode(" \u2190 Kembali ke Cari Jadwal ")
            ];
          }
        }),
        _: 1
      }, _parent));
      if (unref(schedule)) {
        _push(`<div class="schedule-summary" data-v-8c23d228><span class="operator-tag" data-v-8c23d228>${ssrInterpolate(unref(schedule).operator)}</span><span class="route-tag" data-v-8c23d228>${ssrInterpolate(unref(schedule).origin)} \u2794 ${ssrInterpolate(unref(schedule).destination)}</span><span class="date-tag" data-v-8c23d228>\u{1F4C5} ${ssrInterpolate(unref(schedule).departure_date)} (${ssrInterpolate(unref(schedule).departure_time.substring(0, 5))})</span></div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`</div>`);
      if (unref(bookingStore).activeLockedSeat) {
        _push(`<div class="lock-banner-wrapper" data-v-8c23d228>`);
        _push(ssrRenderComponent(_component_CountdownTimer, {
          seconds: unref(bookingStore).lockRemainingSeconds
        }, null, _parent));
        _push(`</div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`<div class="selection-grid" data-v-8c23d228><div class="grid-column" data-v-8c23d228>`);
      _push(ssrRenderComponent(_component_SeatGrid, {
        seats: unref(bookingStore).seats,
        "selected-seat-id": ((_a = unref(selectedSeat)) == null ? void 0 : _a.id) || null,
        onSelectSeat: handleSelectSeat
      }, null, _parent));
      _push(`</div><div class="sidebar-column" data-v-8c23d228><div class="booking-card glass-card" data-v-8c23d228><h3 data-v-8c23d228>Ringkasan Pemesanan</h3>`);
      if (unref(selectedSeat)) {
        _push(`<div class="selected-seat-box" data-v-8c23d228><div class="box-header" data-v-8c23d228><span data-v-8c23d228>Nomor Kursi</span><span class="seat-pill" data-v-8c23d228>${ssrInterpolate(unref(selectedSeat).seat_number)}</span></div><div class="detail-row" data-v-8c23d228><span data-v-8c23d228>Harga Tiket</span><span class="price" data-v-8c23d228>Rp ${ssrInterpolate((_b = unref(schedule)) == null ? void 0 : _b.price.toLocaleString("id-ID"))}</span></div><div class="detail-row" data-v-8c23d228><span data-v-8c23d228>Status Kunci</span>`);
        if (unref(isLockedByMe)) {
          _push(`<span class="badge badge-warning" data-v-8c23d228>Terkunci (5 Menit)</span>`);
        } else {
          _push(`<span class="badge badge-success" data-v-8c23d228>Siap Dikunci</span>`);
        }
        _push(`</div></div>`);
      } else {
        _push(`<div class="no-seat-box" data-v-8c23d228><span class="icon" data-v-8c23d228>\u{1F4BA}</span><p data-v-8c23d228>Silakan klik salah satu kursi hijau yang tersedia pada layout bus di samping.</p></div>`);
      }
      if (unref(message)) {
        _push(`<div class="${ssrRenderClass(["msg-banner", unref(messageType)])}" data-v-8c23d228>${ssrInterpolate(unref(message))}</div>`);
      } else {
        _push(`<!---->`);
      }
      if (unref(selectedSeat) && !unref(isLockedByMe)) {
        _push(`<button class="btn-primary btn-block"${ssrIncludeBooleanAttr(unref(actionLoading)) ? " disabled" : ""} data-v-8c23d228>`);
        if (unref(actionLoading)) {
          _push(`<span data-v-8c23d228>Mengunci...</span>`);
        } else {
          _push(`<span data-v-8c23d228>\u{1F512} Kunci Kursi (5 Menit)</span>`);
        }
        _push(`</button>`);
      } else {
        _push(`<!---->`);
      }
      if (unref(selectedSeat) && unref(isLockedByMe)) {
        _push(`<button class="btn-primary btn-block btn-confirm"${ssrIncludeBooleanAttr(unref(actionLoading)) ? " disabled" : ""} data-v-8c23d228>`);
        if (unref(actionLoading)) {
          _push(`<span data-v-8c23d228>Memproses...</span>`);
        } else {
          _push(`<span data-v-8c23d228>\u2705 Konfirmasi &amp; Bayar Pesanan</span>`);
        }
        _push(`</button>`);
      } else {
        _push(`<!---->`);
      }
      _push(`</div></div></div>`);
      if (unref(confirmedBooking)) {
        _push(`<div class="modal-overlay" data-v-8c23d228><div class="modal-card glass-card" data-v-8c23d228><div class="success-icon" data-v-8c23d228>\u{1F389}</div><h2 data-v-8c23d228>Pemesanan Tiket Berhasil!</h2><p class="modal-sub" data-v-8c23d228>Kode Booking Resmi Anda:</p><div class="ticket-code" data-v-8c23d228>${ssrInterpolate(unref(confirmedBooking).booking_code)}</div><div class="ticket-details" data-v-8c23d228><div class="t-row" data-v-8c23d228><span data-v-8c23d228>Operator:</span> <strong data-v-8c23d228>${ssrInterpolate(unref(confirmedBooking).operator)}</strong></div><div class="t-row" data-v-8c23d228><span data-v-8c23d228>Rute:</span> <strong data-v-8c23d228>${ssrInterpolate(unref(confirmedBooking).origin)} \u2794 ${ssrInterpolate(unref(confirmedBooking).destination)}</strong></div><div class="t-row" data-v-8c23d228><span data-v-8c23d228>Kursi:</span> <strong data-v-8c23d228>${ssrInterpolate(unref(confirmedBooking).seat_number)}</strong></div><div class="t-row" data-v-8c23d228><span data-v-8c23d228>Waktu:</span> <strong data-v-8c23d228>${ssrInterpolate(unref(confirmedBooking).departure_date)} @ ${ssrInterpolate(unref(confirmedBooking).departure_time)}</strong></div><div class="t-row" data-v-8c23d228><span data-v-8c23d228>Total:</span> <strong data-v-8c23d228>Rp ${ssrInterpolate(unref(confirmedBooking).price.toLocaleString("id-ID"))}</strong></div></div><button class="btn-primary btn-block" data-v-8c23d228> Selesai &amp; Kembali ke Beranda </button></div></div>`);
      } else {
        _push(`<!---->`);
      }
      _push(`</div>`);
    };
  }
});
const _sfc_setup = _sfc_main.setup;
_sfc_main.setup = (props, ctx) => {
  const ssrContext = useSSRContext();
  (ssrContext.modules || (ssrContext.modules = /* @__PURE__ */ new Set())).add("pages/schedules/[id].vue");
  return _sfc_setup ? _sfc_setup(props, ctx) : void 0;
};
const _id_ = /* @__PURE__ */ _export_sfc(_sfc_main, [["__scopeId", "data-v-8c23d228"]]);

export { _id_ as default };
//# sourceMappingURL=_id_-CSGordTB.mjs.map
