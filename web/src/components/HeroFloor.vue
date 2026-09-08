<script setup>
import Infographic from "./Infographic.vue";

const holdFig = {
  kind: "hold",
  ledger: 100000,
  holdAmt: 20000,
  ledgerLabel: "остаток в журнале",
  holdLabel: "резерв кассы, проводки нет",
  availNote: "то, что клиент может потратить",
};
</script>

<template>
  <div class="hero-product">
    <Infographic
      kicker="холд"
      plate="01"
      title="Кошелёк после холда"
      :fig="holdFig"
      caption="Касса взяла 20 000. Клиент видит 80 000. В журнале всё ещё 100 000 — деньги не списаны."
    />
  </div>
</template>
