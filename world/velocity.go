package world

type VelocityComponent struct {
	VelocityX float64
	VelocityY float64
}

func (w *World) applyVelocityToEntities() {
	for entity := range w.Entites {
		if !entity.HatVelocityComponent {
			continue
		}

		entity.Position.X += entity.VelocityComponent.VelocityX
		entity.Position.Y += entity.VelocityComponent.VelocityY
	}
}
