<script lang="ts">
  export let data;
</script>

<h1>{data.content.title}</h1>

{#each data.content.sections as section}
  <section>
    <h2>{section.heading}</h2>
    <p>{section.content}</p>
  </section>
{/each}
