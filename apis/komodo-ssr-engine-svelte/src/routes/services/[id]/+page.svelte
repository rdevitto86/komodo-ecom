<script lang="ts">
  export let data;
</script>

<svelte:head>
  <title>{data.meta.title}</title>
  <meta name="description" content={data.meta.description} />
</svelte:head>

<article>
  <h1>{data.content.title}</h1>
  
  {#each data.content.sections as section}
    <section>
      <h2>{section.heading}</h2>
      <div>{@html section.content}</div>
      {#if section.image}
        <img src={section.image} alt={section.heading} />
      {/if}
    </section>
  {/each}
</article>
