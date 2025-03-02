function(ctx) {
  id: ctx.identity.id,
  traits: {
    email: ctx.identity.traits.email,
    name: ctx.identity.traits.name,
    avatar: ctx.identity.traits.avatar,
  },
}
